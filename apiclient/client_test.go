package apiclient

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var accessKeyId string
var secretAccessKey string
var host string

func init() {
	if accessKeyId = os.Getenv("TURBOT_ACCESS_KEY_ID"); accessKeyId == "" {
		log.Fatal("TURBOT_ACCESS_KEY_ID must be set for acceptance tests")
	}
	if secretAccessKey = os.Getenv("TURBOT_SECRET_ACCESS_KEY"); secretAccessKey == "" {
		log.Fatal("TURBOT_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
	if host = os.Getenv("TURBOT_HOST"); host == "" {
		log.Fatal("TURBOT_HOST must be set for acceptance tests")
	}
}

func TestValidateBadCredentials(t *testing.T) {
	client := CreateClient("invalid", "invalid", host)

	err := client.Validate()
	assert.Equal(t, "authorisation failed - have access_key_id and secret_access_key been set correctly?", err.Error())
}

func TestValidateBadHost(t *testing.T) {
	client := CreateClient(accessKeyId, secretAccessKey, "https://bananaman-turbot.invalid.turbot.io")

	err := client.Validate()
	assert.Equal(t, "Post https://bananaman-turbot.invalid.turbot.io: dial tcp: lookup bananaman-turbot.invalid.turbot.io: no such host", err.Error())
}

func TestValidatePass(t *testing.T) {
	client := CreateClient(
		accessKeyId,
		secretAccessKey,
		host)

	err := client.Validate()
	assert.Equal(t, nil, err)
}

func TestLoadSetting_404(t *testing.T) {
	client := CreateClient(
		accessKeyId,
		secretAccessKey,
		host)

	_, err := client.ReadPolicySetting("tmod:@turbot/ssl-check#/policy/types/sslCheck", "166474650907224", "")
	assert.Equal(t, nil, err)

	assert.Equal(t, nil, err)
}
