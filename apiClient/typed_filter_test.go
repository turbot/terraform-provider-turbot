package apiClient

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/assert"
)

func typedFilterClient(t *testing.T, respond func(query string) string) *Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, respond(string(body)))
	}))
	t.Cleanup(server.Close)
	return &Client{Graphql: graphql.NewClient(server.URL + "/graphql")}
}

func TestReadTypedFilterNullIsNotFound(t *testing.T) {
	client := typedFilterClient(t, func(string) string { return `{"data":{"typedFilter":null}}` })

	_, err := client.ReadTypedFilter("1")

	assert.True(t, IsTypedFilterNotFound(err))
}

func TestReadTypedFilterWithoutTypeIsNotFound(t *testing.T) {
	client := typedFilterClient(t, func(string) string {
		return `{"data":{"typedFilter":{"title":null,"filter":null,"type":null,"turbot":{"id":"1"}}}}`
	})

	_, err := client.ReadTypedFilter("1")

	assert.True(t, IsTypedFilterNotFound(err))
}

func TestReadTypedFilterOtherErrorIsNotNotFound(t *testing.T) {
	client := typedFilterClient(t, func(string) string {
		return `{"errors":[{"message":"Forbidden: Insufficient permissions"}]}`
	})

	_, err := client.ReadTypedFilter("1")

	assert.Error(t, err)
	assert.False(t, IsTypedFilterNotFound(err))
}

func TestReadTypedFilterReturnsFilter(t *testing.T) {
	client := typedFilterClient(t, func(string) string {
		return `{"data":{"typedFilter":{"title":"t","description":null,"filter":"title:x","type":"Account","turbot":{"id":"1","akas":null,"tags":{}}}}}`
	})

	filter, err := client.ReadTypedFilter("1")

	assert.NoError(t, err)
	assert.Equal(t, "Account", filter.Type)
	assert.Equal(t, "title:x", filter.Filter)
}

func TestDeleteTypedFilterOfMissingFilterSucceeds(t *testing.T) {
	client := typedFilterClient(t, func(body string) string {
		if strings.Contains(body, "deleteTypedFilter") {
			return `{"data":{"deleteTypedFilter":null},"errors":[{"message":"Bad Request: Failed to get resourceId from aka 1: Not Found"}]}`
		}
		return `{"data":{"typedFilter":null}}`
	})

	assert.NoError(t, client.DeleteTypedFilter("1"))
}

func TestDeleteTypedFilterRefusalSurfaces(t *testing.T) {
	client := typedFilterClient(t, func(body string) string {
		if strings.Contains(body, "deleteTypedFilter") {
			return `{"data":{"deleteTypedFilter":null},"errors":[{"message":"Bad Request: Filter 1 is in use by guardrails and cannot be deleted"}]}`
		}
		return `{"data":{"typedFilter":{"title":"t","filter":"f","type":"Account","turbot":{"id":"1"}}}}`
	})

	err := client.DeleteTypedFilter("1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "in use by guardrails")
}
