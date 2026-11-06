package apiClient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/assert"
)

func userAgentSeenBy(t *testing.T, client *Client) string {
	var seen string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{}}`))
	}))
	defer server.Close()

	client.Graphql = graphql.NewClient(server.URL)
	var resp map[string]interface{}
	assert.NoError(t, client.doRequest(`{ __typename }`, nil, &resp))
	return seen
}

func TestDoRequestSendsConfiguredUserAgent(t *testing.T) {
	ua := "HashiCorp Terraform/1.9.0 (+https://www.terraform.io) terraform-provider-turbot/1.16.0"
	assert.Equal(t, ua, userAgentSeenBy(t, &Client{UserAgent: ua}))
}

func TestDoRequestWithoutUserAgentKeepsGoDefault(t *testing.T) {
	assert.Contains(t, userAgentSeenBy(t, &Client{}), "Go-http-client")
}

func TestCreateClientCarriesUserAgent(t *testing.T) {
	creds := ClientCredentials{AccessKey: "a", SecretKey: "s", Workspace: "https://x.example.com"}
	client, err := CreateClient(ClientConfig{Credentials: creds, UserAgent: "x/1"})
	assert.NoError(t, err)
	assert.Equal(t, "x/1", client.UserAgent)
}
