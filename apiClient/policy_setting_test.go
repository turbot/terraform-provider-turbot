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

// A secret setting must not be silently read through the fallback: the plain fields carry the
// secret REFERENCE {"secret": {"id": "..."}} rather than the value, and storing that would leave a
// permanent diff Terraform would try to "correct". The error must say what grant is missing.
//
// The `secret reference, no type metadata` case is the one that matters most, and is why the guard
// cannot rest on the policy type fields: captured live, PolicyType.secretLevel is null on every
// policy type on the workspace (5000 of 5000) and PolicyType.secret is null for all but a handful,
// so a guard keyed on the metadata alone would let a secret through whenever that metadata is
// absent. The returned value shape is the reliable signal; the metadata is a second one.
func TestReadPolicySettingRefusesSecretOnFallback(t *testing.T) {
	const secretRef = `{"secret":{"id":"387519256421702"}}`
	for _, tc := range []struct {
		name string
		body string
	}{
		{
			"secret reference in value, no type metadata (both fields null)",
			`{"type":{"uri":"tmod:@turbot/azure#/policy/types/clientKey","secret":null,"secretLevel":null},
			  "value":` + secretRef + `,"valueSource":` + secretRef + `,
			  "turbot":{"id":"1","resourceId":"394355648135523"}}`,
		},
		{
			"secret reference in valueSource only",
			`{"type":{"uri":"tmod:@turbot/azure#/policy/types/clientKey"},
			  "value":null,"valueSource":` + secretRef + `,
			  "turbot":{"id":"1","resourceId":"394355648135523"}}`,
		},
		{
			"type marked secret, value withheld entirely",
			`{"type":{"uri":"tmod:@turbot/azure#/policy/types/clientKey","secret":true,"secretLevel":null},
			  "turbot":{"id":"1","resourceId":"394355648135523"}}`,
		},
		{
			"secretLevel CONFIDENTIAL",
			`{"type":{"uri":"tmod:@turbot/aws#/policy/types/secretKey","secret":false,"secretLevel":"CONFIDENTIAL"},
			  "turbot":{"id":"1","resourceId":"394355648135523"}}`,
		},
		{
			"secretLevel SECRET",
			`{"type":{"uri":"tmod:@turbot/aws#/policy/types/secretKey","secret":false,"secretLevel":"SECRET"},
			  "turbot":{"id":"1","resourceId":"394355648135523"}}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stub := &graphqlStub{forbidden: true, fallbackBody: `{"data":{"policySetting":` + tc.body + `}}`}
			client := stubClient(t, stub)

			_, err := client.ReadPolicySetting("1")

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "secret")
			assert.Contains(t, err.Error(), "Turbot/Admin")
			assert.Contains(t, err.Error(), "394355648135523")
		})
	}
}

// The happy path must not be caught by the guard. A non-secret policy type reports null for both
// secret markers (the norm, captured live), so "no metadata" must NOT by itself mean "secret" —
// only an actual secret reference or an explicit marker does.
func TestReadPolicySettingAllowsNonSecretWithNullMetadata(t *testing.T) {
	stub := &graphqlStub{
		forbidden: true,
		fallbackBody: `{"data":{"policySetting":{
			"type":{"uri":"tmod:@turbot/aws-s3#/policy/types/s3AccountPublicAccessBlockSettings","secret":null,"secretLevel":null},
			"value":["Block Public ACLs"],"valueSource":"- \"Block Public ACLs\"\n","precedence":"REQUIRED",
			"turbot":{"id":"394355651429758","resourceId":"394355648135523"}}}}`,
	}
	client := stubClient(t, stub)

	setting, err := client.ReadPolicySetting("394355651429758")

	assert.NoError(t, err, "a non-secret type reports null metadata and must still be readable")
	assert.Equal(t, "394355651429758", setting.Turbot.Id)
	assert.Equal(t, `- "Block Public ACLs"`+"\n", setting.ValueSource)
}

// A valueSource shape that is neither a string nor the known secret reference is refused rather
// than dropped from state — failing closed on anything not understood.
func TestReadPolicySettingRefusesUnexpectedValueSourceShape(t *testing.T) {
	stub := &graphqlStub{
		forbidden: true,
		fallbackBody: `{"data":{"policySetting":{
			"type":{"uri":"tmod:@turbot/aws-s3#/policy/types/x"},
			"value":"ok","valueSource":{"something":"unexpected","and":"another"},
			"turbot":{"id":"1","resourceId":"394355648135523"}}}}`,
	}
	client := stubClient(t, stub)

	_, err := client.ReadPolicySetting("1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected valueSource shape")
}

func TestIsSecretReference(t *testing.T) {
	assert.True(t, isSecretReference(map[string]interface{}{"secret": map[string]interface{}{"id": "1"}}),
		"the live shape must be recognised")
	// a legitimate policy value that merely contains a secret key alongside others is data, not a
	// reference - the live envelope has exactly one key
	assert.False(t, isSecretReference(map[string]interface{}{"secret": "x", "other": "y"}))
	assert.False(t, isSecretReference("a plain string"))
	assert.False(t, isSecretReference([]interface{}{"Block Public ACLs"}))
	assert.False(t, isSecretReference(nil))
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
			"value":"Check: Enabled","valueSource":"Check: Enabled","precedence":"REQUIRED","default":true,
			"turbot":{"id":"394355651429758"}}]}}}`,
	}
	client := stubClient(t, stub)

	setting, err := client.FindPolicySetting("tmod:@turbot/aws-s3#/policy/types/s3AccountPublicAccessBlock", "394355648135523")

	assert.NoError(t, err)
	assert.Equal(t, "394355651429758", setting.Turbot.Id)
	assert.Len(t, stub.requests, 2)
	assert.NotContains(t, stub.requests[1], "secretValue")
}

// Unlike the read, the find path must NOT refuse a secret type — it only answers "does a setting
// already exist". For a secret type the plain fields carry the secret reference, an object: it
// must decode (PolicySetting.ValueSource is a string, so the tolerant type is required) and Value
// must come back non-nil, because the caller's existence check is `Value != nil`. If this returned
// a nil Value, Create would conclude no setting exists and proceed against one that does.
func TestFindPolicySettingDetectsExistingSecretSetting(t *testing.T) {
	const secretRef = `{"secret":{"id":"387519256421702"}}`
	stub := &graphqlStub{
		forbidden: true,
		fallbackBody: `{"data":{"policySettings":{"items":[{
			"value":` + secretRef + `,"valueSource":` + secretRef + `,"precedence":"REQUIRED","default":true,
			"turbot":{"id":"387519256437064"}}]}}}`,
	}
	client := stubClient(t, stub)

	setting, err := client.FindPolicySetting("tmod:@turbot/azure#/policy/types/clientKey", "387519252604383")

	assert.NoError(t, err, "an object-shaped value must decode, not fail with a JSON error")
	assert.NotNil(t, setting.Value, "an existing secret setting must read as existing")
	assert.Equal(t, "387519256437064", setting.Turbot.Id)
}
