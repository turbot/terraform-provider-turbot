package helpers

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

func BuildHttpErrorMessage(err string) error {
	// err = ["graphql", "server returned a non-200 status code", "Code"]
	errCodeString := strings.TrimSpace(strings.Split(err, ":")[2])
	errCode, _ := strconv.ParseUint(errCodeString, 10, 32)

	if int(errCode) == 502 || int(errCode) == 503 || int(errCode) == 504 {
		err = fmt.Sprintf("The server returned a %s (%s). Please wait a few minutes and try again.", http.StatusText(int(errCode)), errCodeString)
		return errors.New(err)
	}
	err = fmt.Sprintf("The server returned a %s (%s). Please contact Turbot support.", http.StatusText(int(errCode)), errCodeString)
	return errors.New(err)
}
