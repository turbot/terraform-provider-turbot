package apiClient

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestCredentialsPrecedence(t *testing.T) {
	type expected struct {
		result bool
		Creds ClientCredentials
	}
	type test struct {
		name string
		Config   ClientConfig
		expected expected
	}
	var tests = []test{
		{
			"Config Access Keys",
			ClientConfig{
				ClientCredentials{
					"xxbd857-XXXX-XXXX-XXXX-xxxxx039ff1x",
					"36xxb4f-XXXX-XXXX-XXXX-c91f44axx4f6",
					"https://example.com/",
				},
				"",
				"",
			},
			expected{
				true,
				ClientCredentials{
					"xxbd857-XXXX-XXXX-XXXX-xxxxx039ff1x",
					"36xxb4f-XXXX-XXXX-XXXX-c91f44axx4f6",
					"https://example.com/",
				},
			},
		},
		{
			"Profile set in config",
			ClientConfig{
				ClientCredentials{
					"",
					"",
					"",
				},
				"",
				"test",
			},
			expected{
				true,
				ClientCredentials{
					"xxbd857-XXXX-XXXX-XXXX-xxxxx039ff1x",
					"36xxb4f-XXXX-XXXX-XXXX-c91f44axx4f6",
					"https://example.com/",
				},
			},
		},
	}
	for _, test := range tests {
		log.Println(test.name)
		credentials, _ := GetCredentials(test.Config)
		assert.Equal(t, test.expected.result, CredentialsSet(credentials))
		assert.ObjectsAreEqual(test.expected.Creds, credentials)
	}
}
