package apiClient

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

/*
[default]
TURBOT_ACCESS_KEY=<valid key>
TURBOT_SECRET_KEY=<valid secret>
turbot_workspace=<valid workspace>

[invalid-keys]
TURBOT_ACCESS_KEY=invalid
TURBOT_SECRET_KEY=invalid
turbot_workspace=https://bananaman-turbot.putney.turbot.io

[invalid-workspace]
TURBOT_ACCESS_KEY=034840a0-fa46-44fa-a445-954c39dxxxx
TURBOT_SECRET_KEY=860c12c8-e8a5-4617-8771-7885a7cbxxxx
turbot_workspace=https://bananaman-turbot.putney.turbot.io_invalid

[unparseable-workspace]
TURBOT_ACCESS_KEY=034840a0-fa46-44fa-a445-954c39d0xxxx
TURBOT_SECRET_KEY=860c12c8-e8a5-4617-8771-7885a7cbxxxx
turbot_workspace=invalid
*/
func TestValidateBadCredentials(t *testing.T) {
	config := ClientConfig{
		Credentials: ClientCredentials{
			AccessKey: "BAD",
			SecretKey: "BAD",
			Workspace: "bananaman-turbot.putney.turbot.io",
		},
	}
	client, err := CreateClient(config)
	assert.Nil(t, err, "error creating client")
	err = client.Validate()
	assert.NotEmpty(t, err)
	assert.Equal(t, "authorisation failed. Verify workspace, access_key and secret_access_key have been set correctly", err.Error())
}

func TestValidateBadWorkspace(t *testing.T) {
	config := ClientConfig{
		Credentials: ClientCredentials{
			AccessKey: "BAD",
			SecretKey: "BAD",
			Workspace: "https://BAD",
		},
	}
	client, err := CreateClient(config)
	assert.Nil(t, err, "error creating client")
	err = client.Validate()
	workspace := config.Credentials.Workspace
	workspaceShort := strings.TrimPrefix(workspace, "https://")
	expected := fmt.Sprintf("Post %s/api/latest/graphql: dial tcp: lookup %s: no such host",
		workspace, workspaceShort)
	assert.NotEmpty(t, err)
	assert.Equal(t, expected, err.Error())
}

func TestValidatePass(t *testing.T) {
	config := ClientConfig{
		Credentials: ClientCredentials{},
	}
	client, err := CreateClient(config)
	assert.Nil(t, err, "error creating client")
	err = client.Validate()
	assert.Equal(t, nil, err)
}

func TestBuildApiUrl(t *testing.T) {
	type urlTest struct {
		url         string
		expectedUrl string
		valid       bool
	}

	tests := []urlTest{
		{url: "https://bananaman-turbot.putney.turbot.io/api/latest", expectedUrl: "https://bananaman-turbot.putney.turbot.io/api/latest/graphql", valid: true},
		{url: "bananaman-turbot.putney.turbot.io/api/v5", expectedUrl: "https://bananaman-turbot.putney.turbot.io/api/v5/graphql", valid: true},
		{url: "bananaman-turbot.putney.turbot.io", expectedUrl: "https://bananaman-turbot.putney.turbot.io/api/latest/graphql", valid: true},
		{url: "bananaman-turbot.putney.turbot.io/", expectedUrl: "https://bananaman-turbot.putney.turbot.io/api/latest/graphql", valid: true},
		{url: "bananaman-turbot.putney.turbot.io/api/latest", expectedUrl: "https://bananaman-turbot.putney.turbot.io/api/latest/graphql", valid: true},
		{url: "bananaman-turbot.putney.turbot.io/api/latest/", expectedUrl: "https://bananaman-turbot.putney.turbot.io/api/latest/graphql", valid: true},
		{url: "bananaman-turbot.putney.turbot.io/api/v5/", expectedUrl: "https://bananaman-turbot.putney.turbot.io/api/v5/graphql", valid: true},
		{url: "https://bananaman-turbot.putney.turbot.io/api/latest/", expectedUrl: "https://bananaman-turbot.putney.turbot.io/api/latest/graphql", valid: true},
		{url: "https://bananaman-turbot.putney.turbot.io/api/v5", expectedUrl: "https://bananaman-turbot.putney.turbot.io/api/v5/graphql", valid: true},
		{url: "https://bananaman-turbot.putney.turbot.io/api/v5/", expectedUrl: "https://bananaman-turbot.putney.turbot.io/api/v5/graphql", valid: true},
	}
	for _, test := range tests {
		url, err := BuildApiUrl(test.url)
		if !test.valid {
			assert.NotEmpty(t, err)
		} else {
			assert.Equal(t, test.expectedUrl, url)
		}

	}
}
