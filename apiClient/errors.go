package apiClient

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func NotFoundError(err error) bool {
	notFoundErr := "(?i)not Found"
	expectedErr := regexp.MustCompile(notFoundErr)
	return expectedErr.Match([]byte(err.Error()))
}

func FailedValidationError(err error) bool {
	dataValidationError := "(?i)data validation failed"
	expectedErr := regexp.MustCompile(dataValidationError)
	return expectedErr.Match([]byte(err.Error()))
}

func BuildHttpErrorMessage(err string) error {
	// if it's a Not Found error, we return the actual graphql error.
	if NotFoundError(errors.New(err)) {
		// if the error is not a 'not found' error, the mod is already installed
		return errors.New(err)
	}
	// err = ["graphql", "server returned a non-200 status code", "Code"]
	errCodeString := strings.TrimSpace(strings.Split(err, ":")[2])
	errCode, _ := strconv.ParseUint(errCodeString, 10, 32)

	// check to see if we fetch correct error code or not
	// return Turbot graphql client error string
	if http.StatusText(int(errCode)) == "" {
		return errors.New(err)
	}
	// retryable error codes - [502, 503, 504]
	if int(errCode) == 502 || int(errCode) == 503 || int(errCode) == 504 {
		err = fmt.Sprintf("The server returned a %s error (%s). Please wait a few minutes and try again.", http.StatusText(int(errCode)), errCodeString)
		return errors.New(err)
	}
	// non-retryable errors
	err = fmt.Sprintf("The server returned a %s error (%s). Please contact Turbot support.", http.StatusText(int(errCode)), errCodeString)
	return errors.New(err)
}
