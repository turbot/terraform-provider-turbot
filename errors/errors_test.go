package errors

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"log"
	"regexp"
	"testing"
)

func TestExtractErrorCode(t *testing.T) {
	type expected struct {
		code int
		err  error
	}
	type test struct {
		name     string
		err      string
		expected expected
	}
	var tests = []test{
		{
			"Bad gateway",
			"graphql: server returned a non-200 status code: 503",
			expected{
				503,
				nil,
			},
		},
		{
			"Not found",
			"graphql:Not Found: Not found error for rocketeer_turbot.grants ",
			expected{
				0,
				errors.Errorf("graphql:Not Found: Not found error for rocketeer_turbot.grants "),
			},
		},
		{
			"Permission denied",
			"graphql: Permission Denied: Insufficient Permissions for rocketeer_turbot.grants  ",
			expected{
				0,
				errors.Errorf("graphql: Permission Denied: Insufficient Permissions for rocketeer_turbot.grants  "),
			},
		},
		{
			"Status network authentication required",
			"graphql: server returned a non-200 status code: 511",
			expected{
				511,
				nil,
			},
		},
		{
			"gRPC error",
			"rpc error: code = Unavailable desc = transport is closing",
			expected{
				0,
				errors.Errorf("rpc error: code = Unavailable desc = transport is closing"),
			},
		},
		{
			"System error",
			"Index out of bound",
			expected{
				0,
				errors.Errorf("Index out of bound"),
			},
		},
		{
			"Bad formatting",
			"graphql: server returned a non-200 status code:       511      ",
			expected{
				511,
				errors.Errorf("Index out of bound"),
			},
		},
	}
	for _, test := range tests {
		log.Println(test.name)
		errCode, err := ExtractErrorCode(errors.Errorf(test.err))
		assert.Equal(t, test.expected.code, errCode)
		assert.ObjectsAreEqual(test.expected.err, err)
	}
}

// The strings below were captured verbatim from a live Guardrails workspace (the punisher
// workspace) while validating Defect 4, plus the provider's own wrappings of them. NotFoundError
// MUST recognise every one, or a genuinely-deleted resource would fail to drop from state.
func TestNotFoundErrorMatchesGenuineNotFound(t *testing.T) {
	genuine := []string{
		// raw Guardrails errors, as returned by the API for a missing item
		"graphql: Not Found: Resource not found or not accessible",
		`graphql: Not Found: Table "punisher_turbot.policy_settings" column "id" with value "999999999999999"`,
		`graphql: Not Found: Table "punisher_turbot.grants" column "id" with value "999999999999999"`,
		`graphql: Not Found: "policy_values" query results. Values: ["173249879813121",["tmod:@turbot/nonexistent#/policy/types/nope"]]`,
		"graphql:Not Found: Not found error for rocketeer_turbot.grants ",
		// provider wrappings (handleReadError / handleUpdateError / handleCreateError)
		"error reading resource: resource not found: 173249879813121",
		"error updating policy pack: resource not found: 12345",
		"error creating smart folder attachment: parent resource not found: 999",
		// lower-case variants, since the match is case-insensitive
		"not found: whatever",
	}
	for _, s := range genuine {
		assert.True(t, NotFoundError(errors.New(s)), "must be recognised as not-found: %q", s)
	}
}

// The core of the defect: an error that merely mentions a missing DEPENDENCY, or any other error
// that happens to contain the word, must NOT be read as "the resource is gone". Several of these
// are strings the provider itself produces.
func TestNotFoundErrorRejectsUnrelated(t *testing.T) {
	unrelated := []string{
		// real provider string, resource_turbot_policy_setting.go — the artifact's example, and
		// the regression this fix exists to prevent
		"policy type tmod:@turbot/aws#/policy/types/foo not found. Is the mod installed?",
		// other real Guardrails / provider errors that are not a missing resource
		`graphql: Cannot query field "thisFieldDoesNotExist" on type "TurbotResourceMetadata".`,
		"graphql: Permission Denied: Insufficient Permissions for rocketeer_turbot.grants",
		"graphql: Data Validation Failed: value is not valid",
		"A policy setting for policy type: 'x', resource: 'y' already exists",
		"the mod could not be found in the registry", // "found" but not the not-found verdict
		"connection refused",
	}
	for _, s := range unrelated {
		assert.False(t, NotFoundError(errors.New(s)), "must NOT be read as a missing resource: %q", s)
	}
}

func TestNotFoundErrorNilIsSafe(t *testing.T) {
	assert.False(t, NotFoundError(nil), "a nil error is not a not-found")
}

// The new matcher must be a STRICT SUBSET of the old `(?i)not found`, so it can only ever match
// fewer errors. That proves it introduces no new false POSITIVES (it never newly classifies an
// unrelated error as not-found). It does NOT prove the absence of new false negatives — matching
// fewer strings is exactly what could miss a genuine not-found shape; that half rests on the
// captured `genuine` corpus, not on this property. Prove the subset relation over a corpus that
// mixes genuine not-founds, unrelated errors, and adversarial near-misses.
func TestNotFoundErrorIsStrictSubsetOfOldMatcher(t *testing.T) {
	old := regexp.MustCompile(`(?i)not found`) // the pre-fix behaviour
	corpus := []string{
		"graphql: Not Found: Resource not found or not accessible",
		"resource not found: 12345",
		"parent resource not found: 999",
		"policy type foo not found. Is the mod installed?",
		"NOT FOUND:",
		"not found",
		"Not Found",
		"nothing was found here",
		"found it",
		"Cannot query field",
		"Permission Denied",
		"",
		"NotFound", // no space — matches neither
		"not_found",
	}
	for _, s := range corpus {
		if NotFoundError(errors.New(s)) {
			assert.True(t, old.MatchString(s),
				"new matcher accepted %q that the old one rejected — not a strict subset", s)
		}
	}
}
