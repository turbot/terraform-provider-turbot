package apiclient

import "regexp"

func NotFoundError(err error) bool {
	notFoundErr := "Not Found"
	expectedErr := regexp.MustCompile(notFoundErr)
	return expectedErr.Match([]byte(err.Error()))
}

func FailedValidationError(err error) bool {
	notFoundErr := "Data validation failed"
	expectedErr := regexp.MustCompile(notFoundErr)
	return expectedErr.Match([]byte(err.Error()))
}
