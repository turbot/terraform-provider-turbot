package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// notFoundRegex matches the two shapes a genuinely-missing resource actually produces:
//
//   - "Not Found:" — the structured prefix Guardrails returns for a missing item, verified live
//     against the workspace ("Not Found: Resource not found or not accessible",
//     `Not Found: Table "..." column "id" with value "..."`, etc.).
//   - "resource not found" — the provider's own wrapping of that error in handleReadError /
//     handleUpdateError / handleCreateError ("error reading X: resource not found: <id>").
//
// Both are matched case-insensitively. The colon in "Not Found:" and the word "resource" are the
// load-bearing anchors: they keep an incidental mention of some OTHER missing thing from being read
// as "the resource is gone". The provider itself emits such strings — e.g.
// resource_turbot_policy_setting.go returns "policy type %s not found. Is the mod installed?" — and
// the previous matcher, `(?i)not found` anywhere, would have treated that as a missing resource and
// dropped a live resource from Terraform state.
//
// This pattern is a strict subset of the old one: every string it matches contains "not found", so
// it can only ever match FEWER errors, never more. That guarantees it introduces no new false
// negatives for the genuine not-found errors callers already rely on.
//
// The complete fix for state membership is to decide from response shape rather than error text
// (see apiClient.ErrTargetNotFound and the `notFound: RETURN_NULL` reads). This narrows the
// text-based fallback that the many callers still using it depend on.
var notFoundRegex = regexp.MustCompile(`(?i)not found:|resource not found`)

func NotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return notFoundRegex.MatchString(err.Error())
}

func FailedValidationError(err error) bool {
	dataValidationError := "(?i)data validation failed"
	expectedErr := regexp.MustCompile(dataValidationError)
	return expectedErr.Match([]byte(err.Error()))
}

func ExtractErrorCode(err error) (int, error) {
	// error returned from machinebox/graphql is of graphql type
	// errorNon200Template = "graphql: server returned a non-200 status code: 503"
	rootError := err
	if strings.Contains(err.Error(), "graphql") {
		errorStringArray := strings.Split(err.Error(), ":")
		if len(errorStringArray) == 3 {
			errCodeString := strings.TrimSpace(errorStringArray[2])
			errCode, err := strconv.ParseUint(errCodeString, 10, 32)
			if err != nil {
				return 0, rootError
			}
			return int(errCode), nil
		}
	}
	return 0, rootError
}

func BuildErrorMessage(err error) error {
	// if it's a Not Found error, we return the actual graphql error.
	if NotFoundError(err) {
		return err
	}
	errCode, err := ExtractErrorCode(err)
	// if we fail to decode the error code, just return the error directly
	if http.StatusText(errCode) == "" {
		return err
	}
	var errString string
	if int(errCode) == 502 || int(errCode) == 503 || int(errCode) == 504 {
		// retryable error codes - [502, 503, 504]
		errString = fmt.Sprintf("The server returned a %s error (%v). Please wait a few minutes and try again.", http.StatusText(errCode), errCode)
	} else {
		// non-retryable errors
		errString = fmt.Sprintf("The server returned a %s error (%v). Please contact Turbot support.", http.StatusText(errCode), errCode)
	}
	return errors.New(errString)
}
