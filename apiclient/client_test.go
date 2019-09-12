package apiclient

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var accessKeyId string
var secretAccessKey string
var workspace string

func init() {
	if accessKeyId = os.Getenv("TURBOT_ACCESS_KEY_ID"); accessKeyId == "" {
		log.Fatal("TURBOT_ACCESS_KEY_ID must be set for client tests")
	}
	if secretAccessKey = os.Getenv("TURBOT_SECRET_ACCESS_KEY"); secretAccessKey == "" {
		log.Fatal("TURBOT_SECRET_ACCESS_KEY must be set for client tests")
	}
	if workspace = os.Getenv("TURBOT_WORKSPACE"); workspace == "" {
		log.Fatal("TURBOT_WORKSPACE must be set for client tests")
	}
}

func TestValidateBadCredentials(t *testing.T) {
	client, err := CreateClient("invalid", "invalid", workspace)
	assert.Nil(t, err, "error creating client")
	err = client.Validate()
	//assert.Equal(t, "graphql: server returned a non-200 status code: 404", err.Error())
	assert.Equal(t, "authorisation failed - have access_key_id and secret_access_key been set correctly?", err.Error())
}

func TestValidateBadWorkspace(t *testing.T) {
	client, err := CreateClient(accessKeyId, secretAccessKey, workspace+"_invalid")
	assert.Nil(t, err, "error creating client")
	err = client.Validate()
	assert.Equal(t, "Post https://bananaman-turbot.putney.turbot.io_invalid/api/v5/graphql: dial tcp: lookup bananaman-turbot.putney.turbot.io_invalid: no such host", err.Error())
}

func TestValidateUnparseableWorkspace(t *testing.T) {
	_, err := CreateClient(accessKeyId, secretAccessKey, "invalid")
	assert.Equal(t, "failed to create client - could not parse workspace url 'invalid'", err.Error())
}

func TestValidatePass(t *testing.T) {
	client, err := CreateClient(accessKeyId, secretAccessKey, workspace)
	assert.Nil(t, err, "error creating client")
	err = client.Validate()
	assert.Equal(t, nil, err)
}
