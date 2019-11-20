package apiClient

import "regexp"

func NotFoundError(err error) bool {
	notFoundErr := "Not Found"
	expectedErr := regexp.MustCompile(notFoundErr)
	return expectedErr.Match([]byte(err.Error()))
}

func FailedValidationError(err error) bool {
	dataValidationError := "data validation failed"
	expectedErr := regexp.MustCompile(dataValidationError)
	return expectedErr.Match([]byte(err.Error()))
}
