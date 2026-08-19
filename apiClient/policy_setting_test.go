package apiClient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/assert"
)

// Guardrails guards the secretValue/secretValueSource FIELDS with a Turbot/Admin check at the
// setting's target resource, even for policy types that hold no secret (verified live: an
// identity with Operator at the Turbot root reads every plain field of a pack-targeted setting,
// but the same query with `value: secretValue` returns "Forbidden: Insufficient permissions for
// resource <packId>"). The without-secrets fallback exists so that identities below Admin can
// still refresh non-secret settings. These tests pin the structural properties of the two
// documents and the fallback decision logic in ReadPolicySetting / FindPolicySetting.

func TestReadPolicySettingQueryShapes(t *testing.T) {
	primary := readPolicySettingQuery()
	fallback := readPolicySettingWithoutSecretsQuery()

	// The primary read keeps requesting the decrypted fields — Admin identities must see real
	// values for secret policies, or drift detection breaks for them.
	assert.Contains(t, primary, "value: secretValue")
	assert.Contains(t, primary, "valueSource: secretValueSource")

	// The whole point of the fallback: no mention of the guarded fields at all. Merely naming
	// secretValue in the document triggers the Admin check, whatever the policy contains.
	assert.NotContains(t, fallback, "secretValue")
	assert.NotContains(t, fallback, "secretValueSource")
	assert.Contains(t, fallback, "value")
	assert.Contains(t, fallback, "valueSource")

	// The fallback must fetch the type's secret markers, so a genuinely secret setting is
	// refused with a clear error instead of silently storing a value the identity cannot
	// decrypt.
	assert.Contains(t, fallback, "secret")
	assert.Contains(t, fallback, "secretLevel")

	// Both documents bind the identifier through a variable, never a literal.
	assert.Contains(t, fallback, "$id: ID!")
}

func TestFindPolicySettingQueryShapes(t *testing.T) {
	fallback := findPolicySettingWithoutSecretsQuery()

	assert.NotContains(t, fallback, "secretValue")
	assert.NotContains(t, fallback, "secretValueSource")
	assert.Contains(t, fallback, "policySettingList(filter: $filter)")
	assert.Contains(t, fallback, "$filter: [String!]")
}

// graphqlStub answers each request by inspecting the query document: documents that mention
// secretValue get the `forbidden` payload, everything else gets `fallbackBody`. It records the
// documents it served so tests can assert how many calls were made and with which shape.
type graphqlStub struct {
	forbidden    bool
	fallbackBody string
	requests     []string
}

func (s *graphqlStub) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Query string `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.requests = append(s.requests, req.Query)
		w.Header().Set("Content-Type", "application/json")
		if s.forbidden && strings.Contains(req.Query, "secretValue") {
			fmt.Fprint(w, `{"errors":[{"message":"Forbidden: Insufficient permissions for resource 394355648135523"}]}`)
			return
		}
		fmt.Fprint(w, s.fallbackBody)
	}
}

func stubClient(t *testing.T, stub *graphqlStub) *Client {
	t.Helper()
	server := httptest.NewServer(stub.handler())
	t.Cleanup(server.Close)
	return &Client{Graphql: graphql.NewClient(server.URL + "/graphql")}
}

func TestReadPolicySettingFallsBackWithoutSecretsOnForbidden(t *testing.T) {
	stub := &graphqlStub{
		forbidden: true,
		fallbackBody: `{"data":{"policySetting":{
			"type":{"uri":"tmod:@turbot/aws-s3#/policy/types/s3AccountPublicAccessBlockSettings","secret":false,"secretLevel":"NONE"},
			"value":["Block Public ACLs"],
			"valueSource":"- \"Block Public ACLs\"\n",
			"precedence":"REQUIRED",
			"turbot":{"id":"394355651429758","resourceId":"394355648135523"}}}}`,
	}
	client := stubClient(t, stub)

	setting, err := client.ReadPolicySetting("394355651429758")

	assert.NoError(t, err)
	assert.Equal(t, "394355651429758", setting.Turbot.Id)
	assert.Equal(t, "REQUIRED", setting.Precedence)
	assert.Equal(t, `- "Block Public ACLs"`+"\n", setting.ValueSource)
	// exactly two calls: the primary (Forbidden) and the without-secrets retry
	assert.Len(t, stub.requests, 2)
	assert.Contains(t, stub.requests[0], "secretValue")
	assert.NotContains(t, stub.requests[1], "secretValue")
}

// A secret setting must not be silently read through the fallback: the plain fields do not carry
// the real value, and storing them would corrupt state. The error must say what grant is missing.
func TestReadPolicySettingRefusesSecretTypeOnFallback(t *testing.T) {
	for _, tc := range []struct {
		name string
		typ  string
	}{
		{"secretLevel CONFIDENTIAL", `{"uri":"tmod:@turbot/aws#/policy/types/secretKey","secret":false,"secretLevel":"CONFIDENTIAL"}`},
		{"secretLevel SECRET", `{"uri":"tmod:@turbot/aws#/policy/types/secretKey","secret":false,"secretLevel":"SECRET"}`},
		{"legacy secret flag", `{"uri":"tmod:@turbot/aws#/policy/types/secretKey","secret":true,"secretLevel":""}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stub := &graphqlStub{
				forbidden: true,
				fallbackBody: `{"data":{"policySetting":{"type":` + tc.typ + `,
					"turbot":{"id":"1","resourceId":"394355648135523"}}}}`,
			}
			client := stubClient(t, stub)

			_, err := client.ReadPolicySetting("1")

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "secret")
			assert.Contains(t, err.Error(), "Turbot/Admin")
			assert.Contains(t, err.Error(), "394355648135523")
		})
	}
}

// If the fallback fails too, the denial is on the setting itself (identity below Metadata, or a
// workspace whose schema predates the fallback's type fields). The ORIGINAL error must surface —
// degrading to exactly the pre-fallback behavior, never a confusing secondary error.
func TestReadPolicySettingSurfacesOriginalErrorWhenFallbackFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"errors":[{"message":"Forbidden: Insufficient permissions for resource 394355648135523"}]}`)
	}))
	defer server.Close()
	client := &Client{Graphql: graphql.NewClient(server.URL + "/graphql")}

	_, err := client.ReadPolicySetting("394355651429758")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading policy setting")
	assert.Contains(t, err.Error(), "Forbidden: Insufficient permissions for resource 394355648135523")
}

// A non-Forbidden error must not trigger the fallback: retrying a Not Found would waste a round
// trip, and NotFoundError classification (which Exists relies on to clear state) must be
// preserved by handleReadError.
func TestReadPolicySettingDoesNotFallBackOnOtherErrors(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"errors":[{"message":"Not Found: Resource not found or not accessible"}]}`)
	}))
	defer server.Close()
	client := &Client{Graphql: graphql.NewClient(server.URL + "/graphql")}

	_, err := client.ReadPolicySetting("394355651429758")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource not found")
	assert.Equal(t, 1, calls, "a non-Forbidden error must not trigger the fallback read")
}

func TestReadPolicySettingSingleCallOnSuccess(t *testing.T) {
	stub := &graphqlStub{
		forbidden: false,
		fallbackBody: `{"data":{"policySetting":{
			"type":{"uri":"tmod:@turbot/aws-s3#/policy/types/s3AccountPublicAccessBlockSettings"},
			"value":"ok","precedence":"REQUIRED",
			"turbot":{"id":"394355651429758","resourceId":"394355648135523"}}}}`,
	}
	client := stubClient(t, stub)

	setting, err := client.ReadPolicySetting("394355651429758")

	assert.NoError(t, err)
	assert.Equal(t, "394355651429758", setting.Turbot.Id)
	assert.Len(t, stub.requests, 1, "no fallback call when the primary read succeeds")
}

// FindPolicySetting feeds duplicate detection in Create. A Forbidden caused by the secretValue
// guard on a matched item must fall back rather than fail, so Create can produce its correct
// "already exists, import it" error instead of a bare Forbidden.
func TestFindPolicySettingFallsBackOnForbidden(t *testing.T) {
	stub := &graphqlStub{
		forbidden: true,
		fallbackBody: `{"data":{"policySettings":{"items":[{
			"value":"Check: Enabled","precedence":"REQUIRED","default":true,
			"turbot":{"id":"394355651429758"}}]}}}`,
	}
	client := stubClient(t, stub)

	setting, err := client.FindPolicySetting("tmod:@turbot/aws-s3#/policy/types/s3AccountPublicAccessBlock", "394355648135523")

	assert.NoError(t, err)
	assert.Equal(t, "394355651429758", setting.Turbot.Id)
	assert.Len(t, stub.requests, 2)
	assert.NotContains(t, stub.requests[1], "secretValue")
}
