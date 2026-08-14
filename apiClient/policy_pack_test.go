package apiClient

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The whole point of readPolicyPackIdentityQuery is that it does NOT use the generic
// `resource(id:)` query, which Guardrails authorizes against the pack itself. Packs live at
// the Turbot root, so `resource(id:)` requires a root-level grant, while attaching a pack
// only requires permissions on the attachment target. Assert the distinction directly so a
// future refactor back to `resource(id:)` fails loudly rather than silently reintroducing
// the permission regression.
func TestReadPolicyPackIdentityQuery(t *testing.T) {
	var tests = []struct {
		name       string
		policyPack string
	}{
		{"numeric id", "343645598782139"},
		{"aka", "aws_s3_bucket_versioning_enabled"},
		{"mod-style aka", "tmod:@turbot/turbot#/my-pack"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			query := readPolicyPackIdentityQuery(test.policyPack)

			assert.Contains(t, query, "policyPack(id:\""+test.policyPack+"\")",
				"query must address the pack through policyPack(id:)")
			assert.NotContains(t, query, "resource(id:",
				"query must NOT use resource(id:) — that requires a grant on the pack itself")
			// both fields are required: id resolves the attach mutation input, akas feed
			// suppressIfAkaMatches so an aka in config does not diff against an id in state
			assert.Contains(t, query, "id")
			assert.Contains(t, query, "akas")
		})
	}
}

// Exists() answers "is this pack attached to this resource" from the resource side. Reading
// the pack instead would need a grant on the pack, and would depend on a pack-wide list that
// can span the whole hierarchy.
func TestReadAttachedPolicyPacksQuery(t *testing.T) {
	var tests = []struct {
		name     string
		resource string
	}{
		{"numeric id", "191926035367605"},
		{"aka", "my_bu_folder"},
		{"arn-style aka", "arn:aws:::111122223333"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			query := readAttachedPolicyPacksQuery(test.resource)

			assert.Contains(t, query, "resource(id:\""+test.resource+"\")",
				"query must be anchored on the attachment target, not the pack")
			assert.Contains(t, query, "attachedSmartFolders",
				"must read the packs attached to the resource")
			assert.NotContains(t, query, "attachedResources",
				"attachedResources is the pack-side list and spans the hierarchy")
			assert.NotContains(t, query, "policyPack(id:",
				"Exists must not need to read the pack at all")
		})
	}
}

// A pack may be identified in state or config by its numeric id or by any of its akas, so
// matching has to accept either. Both turbot_policy_pack_attachment and
// turbot_smart_folder_attachment depend on this.
func TestPolicyPackInList(t *testing.T) {
	attached := []TurbotResourceMetadata{
		{Id: "343645598782139", Akas: []string{"aws_account_decommission", "decomm"}},
		{Id: "325675724117430", Akas: nil},
	}

	var tests = []struct {
		name       string
		list       []TurbotResourceMetadata
		policyPack string
		expected   bool
	}{
		{"matches by numeric id", attached, "343645598782139", true},
		{"matches by first aka", attached, "aws_account_decommission", true},
		{"matches by second aka", attached, "decomm", true},
		{"matches a pack with no akas by id", attached, "325675724117430", true},
		{"does not match an unrelated id", attached, "999999999999999", false},
		{"does not match an unrelated aka", attached, "some_other_pack", false},
		{"empty list matches nothing", []TurbotResourceMetadata{}, "343645598782139", false},
		{"nil list matches nothing", nil, "343645598782139", false},
		{"empty needle does not match a pack with no akas", attached, "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, policyPackInList(test.list, test.policyPack))
		})
	}
}

// The queries are built with fmt.Sprintf, so guard against a caller-supplied identifier
// breaking out of the string literal and corrupting the query shape.
func TestPolicyPackQueriesRemainWellFormed(t *testing.T) {
	for _, id := range []string{"343645598782139", "aws_s3_bucket_versioning_enabled"} {
		identity := readPolicyPackIdentityQuery(id)
		attached := readAttachedPolicyPacksQuery(id)
		for name, query := range map[string]string{"identity": identity, "attached": attached} {
			assert.Equal(t, strings.Count(query, "{"), strings.Count(query, "}"),
				"%s query braces must balance", name)
			assert.Equal(t, 2, strings.Count(query, "\""),
				"%s query must contain exactly one quoted identifier", name)
		}
	}
}
