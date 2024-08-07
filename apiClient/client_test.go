package apiClient

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredentialsPrecedence(t *testing.T) {
	type expected struct {
		result bool
		Creds  ClientCredentials
	}
	type test struct {
		name     string
		Config   ClientConfig
		expected expected
	}
	var tests = []test{
		{
			"Config has credentials",
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
		// TODO: To make the below test work, create `test` profile in config/turbot/credentials.yml file
		//{
		//	"Config has profile",
		//	ClientConfig{
		//		ClientCredentials{
		//			"",
		//			"",
		//			"",
		//		},
		//		"",
		//		"test",
		//	},
		//	expected{
		//		true,
		//		ClientCredentials{
		//			"xxbd857-XXXX-XXXX-XXXX-xxxxx039ff1x",
		//			"36xxb4f-XXXX-XXXX-XXXX-c91f44axx4f6",
		//			"https://example.com/",
		//		},
		//	},
		//},
		{
			"Empty Config",
			ClientConfig{
				ClientCredentials{
					"",
					"",
					"",
				},
				"",
				"",
			},
			expected{
				true,
				ClientCredentials{
					os.Getenv("TURBOT_ACCESS_KEY"),
					os.Getenv("TURBOT_SECRET_KEY"),
					os.Getenv("TURBOT_WORKSPACE"),
				},
			},
		},
	}
	for _, test := range tests {
		log.Println(test.name)
		credentials, _ := GetCredentials(test.Config)
		if !CredentialsSet(credentials) {
			fmt.Printf(`In order to successfully execute TF calls, credentials must me set`)
		} else {
			assert.Equal(t, test.expected.result, CredentialsSet(credentials))
			assert.ObjectsAreEqual(test.expected.Creds, credentials)
		}
	}
}
