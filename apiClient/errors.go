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

func BuildHttpErrorMessage(err error) error {
	// if it's a Not Found error, we return the actual graphql error.
	if NotFoundError(err) {
		return err
	}
	errCodeString := strings.TrimSpace(strings.Split(err.Error(), ":")[2])
	errCode, _ := strconv.ParseUint(errCodeString, 10, 32)

	// if we fail to decode the error code, just return the error directly
	if http.StatusText(int(errCode)) == "" {
		return err
	}
	var errString string
	if int(errCode) == 502 || int(errCode) == 503 || int(errCode) == 504 {
		// retryable error codes - [502, 503, 504]
		errString = fmt.Sprintf("The server returned a %s error (%s). Please wait a few minutes and try again.", http.StatusText(int(errCode)), errCodeString)
	} else {
		// non-retryable errors
		errString = fmt.Sprintf("The server returned a %s error (%s). Please contact Turbot support.", http.StatusText(int(errCode)), errCodeString)
	}
	return errors.New(errString)
}
