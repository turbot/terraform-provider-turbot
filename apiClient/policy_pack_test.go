package apiClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The whole point of readPolicyPackIdentityQuery is that it does NOT use the generic
// `resource(id:)` query, which Guardrails authorizes against the pack itself. Packs live at
// the Turbot root by default, so `resource(id:)` requires a grant wherever the pack sits, while
// attaching a pack
// only requires permissions on the attachment target. Assert the distinction directly so a
// future refactor back to `resource(id:)` fails loudly rather than silently reintroducing
// the permission regression.
func TestReadPolicyPackIdentityQuery(t *testing.T) {
	query := readPolicyPackIdentityQuery()

	assert.Contains(t, query, "policyPack(id: $id)",
		"query must address the pack through policyPack(id:)")
	assert.NotContains(t, query, "resource(id:",
		"query must NOT use resource(id:) - that requires a grant on the pack itself")
	// both fields are required: id resolves the attach mutation input, akas feed
	// suppressIfAkaMatches so an aka in config does not diff against an id in state
	assert.Contains(t, query, "id")
	assert.Contains(t, query, "akas")
	// the identifier must arrive as a variable, never interpolated
	assert.Contains(t, query, "$id: ID!")
}

// Exists() answers "is this pack attached to this resource" from the resource side. Reading
// the pack instead would need a grant on the pack, and would depend on a pack-wide list that
// can span the whole hierarchy.
func TestReadAttachedPolicyPacksQuery(t *testing.T) {
	query := readAttachedPolicyPacksQuery()

	assert.Contains(t, query, "resource(id: $id",
		"query must be anchored on the attachment target, not the pack")
	assert.Contains(t, query, "attachedSmartFolders",
		"must read the packs attached to the resource")
	assert.NotContains(t, query, "attachedResources",
		"attachedResources is the pack-side list and spans the hierarchy")
	assert.NotContains(t, query, "policyPack(id:",
		"Exists must not need to read the pack at all")
	// attachedSmartFolders takes no paging arguments, so paging.next is the only signal
	// that the list was truncated - a truncated list would make Exists false-negative
	assert.Contains(t, query, "paging",
		"must select paging so truncation can be detected")
	assert.Contains(t, query, "next")
	// the identifier must arrive as a variable, never interpolated
	assert.Contains(t, query, "$id: ID!")
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

// Both read builders pass the identifier as a GraphQL variable, so no caller-supplied value can
// reach the query document at all. This is the regression guard for that: previously the builders
// interpolated with fmt.Sprintf, and an identifier containing a double quote escaped the string
// literal and appended attacker-chosen GraphQL, executed with the provider's credentials.
//
// The payload below is the proof-of-concept from review. Against the old interpolating builders it
// produced a valid document with two top-level fields - and it kept braces balanced (5 open / 5
// close), so a brace-balance assertion passed it. Only the queries being identifier-free makes the
// class impossible, which is what this asserts.
func TestReadBuildersCannotBeInjected(t *testing.T) {
	payload := `x") { turbot { id } } deleteMe: resource(id:"12345`

	for name, query := range map[string]string{
		"identity": readPolicyPackIdentityQuery(),
		"attached": readAttachedPolicyPacksQuery(),
	} {
		// the document is a constant: nothing a caller supplies can appear in it
		assert.NotContains(t, query, payload, "%s query must not embed caller input", name)
		assert.NotContains(t, query, "deleteMe", "%s query must not embed injected fields", name)
		assert.NotContains(t, query, "12345", "%s query must not embed caller input", name)

		// exactly one top-level selection, so no second field can be smuggled in
		assert.Equal(t, 1, strings.Count(query, "$id: ID!"), "%s query declares one variable", name)
		assert.Equal(t, strings.Count(query, "{"), strings.Count(query, "}"),
			"%s query braces must balance", name)

		// no string literals at all - the only way to interpolate is to add one
		assert.NotContains(t, query, `"`, "%s query must contain no string literal", name)
	}
}

// Under truncation only ABSENCE is ambiguous, so PolicyPackAttached errors on absent-and-truncated
// and answers true whenever the pack is in the page it did read. The guard's whole correctness
// rests on paging.next decoding onto Paging.Next, so exercise the real JSON mapping rather than a
// struct assignment - a renamed field must fail this test, not slip through it.
func TestAttachedPolicyPacksTruncationIsDetectable(t *testing.T) {
	// A cursor must decode, or the guard can never fire.
	var truncated AttachedPolicyPacksResponse
	err := json.Unmarshal([]byte(`{"resource":{"attachedSmartFolders":{"paging":{"next":"cursor-abc"},"items":[]}}}`), &truncated)
	assert.NoError(t, err)
	assert.Equal(t, "cursor-abc", truncated.Resource.AttachedSmartFolders.Paging.Next,
		"paging.next must decode onto Paging.Next")

	// A complete list sends null, which must land as empty - what the guard reads as complete.
	var complete AttachedPolicyPacksResponse
	err = json.Unmarshal([]byte(`{"resource":{"attachedSmartFolders":{"paging":{"next":null},"items":[]}}}`), &complete)
	assert.NoError(t, err)
	assert.Equal(t, "", complete.Resource.AttachedSmartFolders.Paging.Next)

	// And the items alongside a cursor must still decode, since a pack found in a truncated page
	// is a definitive yes.
	var withItems AttachedPolicyPacksResponse
	err = json.Unmarshal([]byte(`{"resource":{"attachedSmartFolders":{"paging":{"next":"c"},"items":[{"turbot":{"id":"123","akas":["my_pack"]}}]}}}`), &withItems)
	assert.NoError(t, err)
	assert.Len(t, withItems.Resource.AttachedSmartFolders.Items, 1)
	assert.Equal(t, "123", withItems.Resource.AttachedSmartFolders.Items[0].Turbot.Id)
	assert.Equal(t, []string{"my_pack"}, withItems.Resource.AttachedSmartFolders.Items[0].Turbot.Akas)
}

// Exists decides state membership from this, so it must be driven by response SHAPE rather than by
// matching "not found" in error text. errors.NotFoundError matches `(?i)not found` anywhere, so an
// unrelated not-found would drop a live attachment out of state and have Terraform recreate it.
func TestTargetNotFoundIsTypedNotTextMatched(t *testing.T) {
	// The sentinel is recognised through wrapping, which is how ReadAttachedPolicyPacks returns it.
	wrapped := fmt.Errorf("%w: %s", ErrTargetNotFound, "391406345032847")
	assert.True(t, IsTargetNotFound(wrapped), "wrapped sentinel must be recognised")
	assert.Contains(t, wrapped.Error(), "391406345032847", "the identifier stays in the message")

	// Errors that merely mention "not found" must NOT be treated as a missing target - this is the
	// regression this replaces.
	for _, text := range []string{
		"policy type not found",
		"error reading attached policy packs: mod not found",
		"Not Found: something entirely unrelated",
	} {
		assert.False(t, IsTargetNotFound(errors.New(text)),
			"%q must not be read as a missing target", text)
	}

	// And a nil error is not a missing target.
	assert.False(t, IsTargetNotFound(nil))
}

// A `resource: null` reply is the only thing that means "target missing", so the response type must
// be able to represent it - a value struct would decode null as a zero struct with no attachments,
// which is a completely different fact.
func TestAttachedPolicyPacksDecodesNullResource(t *testing.T) {
	var missing AttachedPolicyPacksResponse
	err := json.Unmarshal([]byte(`{"resource":null}`), &missing)
	assert.NoError(t, err)
	assert.Nil(t, missing.Resource, "a null resource must decode as nil, not as an empty struct")

	var present AttachedPolicyPacksResponse
	err = json.Unmarshal([]byte(`{"resource":{"attachedSmartFolders":{"paging":{"next":null},"items":[]}}}`), &present)
	assert.NoError(t, err)
	assert.NotNil(t, present.Resource, "an existing target with no attachments is not nil")
	assert.Empty(t, present.Resource.AttachedSmartFolders.Items)
}

// The query must actually ask for RETURN_NULL, or a missing target arrives as an error and the
// typed check above can never fire.
func TestReadAttachedPolicyPacksQueryRequestsNullOnMissing(t *testing.T) {
	query := readAttachedPolicyPacksQuery()
	assert.Contains(t, query, "options: {notFound: RETURN_NULL}",
		"a missing target must come back as null, not as an error")
}
