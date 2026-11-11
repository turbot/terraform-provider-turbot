package turbot

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/machinebox/graphql"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
)

// test suites
func TestAccSamlDirectory_Basic(t *testing.T) {
	resourceName := "turbot_saml_directory.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSamlDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSamlDirectoryConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSamlDirectoryExists("turbot_saml_directory.test"),
					resource.TestCheckResourceAttr("turbot_saml_directory.test", "title", "provider-test"),
					resource.TestCheckResourceAttr("turbot_saml_directory.test", "profile_id_template", "{{profile.email}}"),
					resource.TestCheckResourceAttr("turbot_saml_directory.test", "description", "SAML Directory Testing1"),
					resource.TestCheckResourceAttr("turbot_saml_directory.test", "entry_point", "https://example.com/myapp/sso/saml"),
					resource.TestCheckResourceAttr("turbot_saml_directory.test", "certificate", "-----BEGIN CERTIFICATE-----\nMIICiTCCAfICCQD6m7oRw0uXOjANBgkqhkiG9w0BAQUFADCBiDELMAkGA1UEBhMC\nVVMxCzAJBgNVBAgTAldBMRAwDgYDVQQHEwdTZWF0dGxlMQ8wDQYDVQQKEwZBbWF6\nb24xFDASBgNVBAsTC0lBTSBDb25zb2xlMRIwEAYDVQQDEwlUZXN0Q2lsYWMxHzAd\nBgkqhkiG9w0BCQEWEG5vb25lQGFtYXpvbi5jb20wHhcNMTEwNDI1MjA0NTIxWhcN\nMTIwNDI0MjA0NTIxWjCBiDELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAldBMRAwDgYD\nVQQHEwdTZWF0dGxlMQ8wDQYDVQQKEwZBbWF6b24xFDASBgNVBAsTC0lBTSBDb25z\nb2xlMRIwEAYDVQQDEwlUZXN0Q2lsYWMxHzAdBgkqhkiG9w0BCQEWEG5vb25lQGFt\nYXpvbi5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMaK0dn+a4GmWIWJ\n21uUSfwfEvySWtC2XADZ4nB+BLYgVIk60CpiwsZ3G93vUEIO3IyNoH/f0wYK8m9T\nrDHudUZg3qX4waLG5M43q7Wgc/MbQITxOUSQv7c7ugFFDzQGBzZswY6786m86gpE\nIbb3OhjZnzcvQAaRHhdlQWIMm2nrAgMBAAEwDQYJKoZIhvcNAQEFBQADgYEAtCu4\nnUhVVxYUntneD9+h8Mg9q6q+auNKyExzyLwaxlAoo7TJHidbtS4J5iNmZgXL0Fkb\nFFBjvSfpJIlJ00zbhNYS5f6GuoEDmFJl0ZxBHjJnyp378OD8uTs7fLvjx79LjSTb\nNYiytVbZPQUQ5Yaxu2jXnimvw3rrszlaEXAMPLE=\n-----END CERTIFICATE-----"),
				),
			},
			{
				Config: testAccSamlDirectoryUpdateDescConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSamlDirectoryExists("turbot_saml_directory.test"),
					resource.TestCheckResourceAttr("turbot_saml_directory.test", "description", "SAML Directory Testing1 updated"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSamlDirectory_EnumValues(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSamlDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSamlDirectoryOptionalAttributeConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSamlDirectoryExists("turbot_saml_directory.test"),
					resource.TestCheckResourceAttr("turbot_saml_directory.test", "title", "provider-test"),
					resource.TestCheckResourceAttr("turbot_saml_directory.test", "entry_point", "https://example.com/myapp/sso/saml"),
					resource.TestCheckResourceAttr("turbot_saml_directory.test", "profile_id_template", "{{profile.email}}"),
					resource.TestCheckResourceAttr("turbot_saml_directory.test", "name_id_format", "UNSPECIFIED"),
					resource.TestCheckResourceAttr("turbot_saml_directory.test", "allow_group_syncing", "true"),
				),
			},
		},
	})
}

func TestAccSamlDirectory_Signing(t *testing.T) {
	resourceName := "turbot_saml_directory.test"
	key, certificate := testAccSamlSigningKeyPair(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSamlDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSamlDirectorySigningConfig(key, certificate, `
	sign_requests       = true
	signature_algorithm = "SHA256"`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSamlDirectoryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sign_requests", "true"),
					resource.TestCheckResourceAttr(resourceName, "signature_algorithm", "SHA256"),
					resource.TestCheckResourceAttr(resourceName, "signature_private_key", key),
				),
			},
			{
				Config: testAccSamlDirectorySigningConfig(key, certificate, `
	sign_requests       = false
	signature_algorithm = "SHA512"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "sign_requests", "false"),
					resource.TestCheckResourceAttr(resourceName, "signature_algorithm", "SHA512"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Guardrails returns only a reference to the stored key.
				ImportStateVerifyIgnore: []string{"signature_private_key"},
			},
		},
	})
}

// A directory that is saved in the console must not change the plan of a configuration that leaves signing unset.
func TestAccSamlDirectory_SavedInConsole(t *testing.T) {
	resourceName := "turbot_saml_directory.test"
	_, certificate := testAccSamlSigningKeyPair(t)
	var id string
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSamlDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSamlDirectorySigningConfig("", certificate, ""),
				Check: func(s *terraform.State) error {
					id = s.RootModule().Resources[resourceName].Primary.ID
					return nil
				},
			},
			{
				// The console's edit form saves signRequests and signatureAlgorithm even when signing is off.
				PreConfig: func() {
					client := testAccProvider.Meta().(*apiClient.Client)
					input := map[string]interface{}{"id": id, "signRequests": false, "signatureAlgorithm": "SHA1"}
					if _, err := client.UpdateSamlDirectory(input); err != nil {
						t.Fatalf("saving as the console does: %s", err)
					}
				},
				Config:   testAccSamlDirectorySigningConfig("", certificate, ""),
				PlanOnly: true,
			},
		},
	})
}

// testAccSamlSigningKeyPair returns a throwaway PEM private key and a self-signed certificate.
func testAccSamlSigningKeyPair(t *testing.T) (string, string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{SerialNumber: big.NewInt(1), NotAfter: time.Now().Add(time.Hour)}
	cert, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})),
		string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert}))
}

func testAccSamlDirectorySigningConfig(key, certificate, signing string) string {
	keyAttribute := ""
	if key != "" {
		keyAttribute = fmt.Sprintf("signature_private_key = <<EOT\n%sEOT\n", key)
	}
	return fmt.Sprintf(`
resource "turbot_saml_directory" "test" {
	title               = "provider-test"
	parent              = "tmod:@turbot/turbot#/"
	profile_id_template = "{{profile.email}}"
	entry_point         = "https://example.com/myapp/sso/saml"
	certificate         = <<EOT
%sEOT
	%s%s
}
`, certificate, keyAttribute, signing)
}

// configs
func testAccSamlDirectoryConfig() string {
	return `
resource "turbot_saml_directory" "test" {
	title               = "provider-test"
	parent              = "tmod:@turbot/turbot#/"
	profile_id_template = "{{profile.email}}"
	description         = "SAML Directory Testing1"
	entry_point         = "https://example.com/myapp/sso/saml"
	certificate         = "-----BEGIN CERTIFICATE-----\nMIICiTCCAfICCQD6m7oRw0uXOjANBgkqhkiG9w0BAQUFADCBiDELMAkGA1UEBhMC\nVVMxCzAJBgNVBAgTAldBMRAwDgYDVQQHEwdTZWF0dGxlMQ8wDQYDVQQKEwZBbWF6\nb24xFDASBgNVBAsTC0lBTSBDb25zb2xlMRIwEAYDVQQDEwlUZXN0Q2lsYWMxHzAd\nBgkqhkiG9w0BCQEWEG5vb25lQGFtYXpvbi5jb20wHhcNMTEwNDI1MjA0NTIxWhcN\nMTIwNDI0MjA0NTIxWjCBiDELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAldBMRAwDgYD\nVQQHEwdTZWF0dGxlMQ8wDQYDVQQKEwZBbWF6b24xFDASBgNVBAsTC0lBTSBDb25z\nb2xlMRIwEAYDVQQDEwlUZXN0Q2lsYWMxHzAdBgkqhkiG9w0BCQEWEG5vb25lQGFt\nYXpvbi5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMaK0dn+a4GmWIWJ\n21uUSfwfEvySWtC2XADZ4nB+BLYgVIk60CpiwsZ3G93vUEIO3IyNoH/f0wYK8m9T\nrDHudUZg3qX4waLG5M43q7Wgc/MbQITxOUSQv7c7ugFFDzQGBzZswY6786m86gpE\nIbb3OhjZnzcvQAaRHhdlQWIMm2nrAgMBAAEwDQYJKoZIhvcNAQEFBQADgYEAtCu4\nnUhVVxYUntneD9+h8Mg9q6q+auNKyExzyLwaxlAoo7TJHidbtS4J5iNmZgXL0Fkb\nFFBjvSfpJIlJ00zbhNYS5f6GuoEDmFJl0ZxBHjJnyp378OD8uTs7fLvjx79LjSTb\nNYiytVbZPQUQ5Yaxu2jXnimvw3rrszlaEXAMPLE=\n-----END CERTIFICATE-----"
}
`
}

func testAccSamlDirectoryUpdateDescConfig() string {
	return `
resource "turbot_saml_directory" "test" {
	title               = "provider-test"
	parent              = "tmod:@turbot/turbot#/"
	profile_id_template = "{{profile.email}}"
	description         = "SAML Directory Testing1 updated"
	entry_point         = "https://example.com/myapp/sso/saml"
	certificate         = "-----BEGIN CERTIFICATE-----\nMIICiTCCAfICCQD6m7oRw0uXOjANBgkqhkiG9w0BAQUFADCBiDELMAkGA1UEBhMC\nVVMxCzAJBgNVBAgTAldBMRAwDgYDVQQHEwdTZWF0dGxlMQ8wDQYDVQQKEwZBbWF6\nb24xFDASBgNVBAsTC0lBTSBDb25zb2xlMRIwEAYDVQQDEwlUZXN0Q2lsYWMxHzAd\nBgkqhkiG9w0BCQEWEG5vb25lQGFtYXpvbi5jb20wHhcNMTEwNDI1MjA0NTIxWhcN\nMTIwNDI0MjA0NTIxWjCBiDELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAldBMRAwDgYD\nVQQHEwdTZWF0dGxlMQ8wDQYDVQQKEwZBbWF6b24xFDASBgNVBAsTC0lBTSBDb25z\nb2xlMRIwEAYDVQQDEwlUZXN0Q2lsYWMxHzAdBgkqhkiG9w0BCQEWEG5vb25lQGFt\nYXpvbi5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMaK0dn+a4GmWIWJ\n21uUSfwfEvySWtC2XADZ4nB+BLYgVIk60CpiwsZ3G93vUEIO3IyNoH/f0wYK8m9T\nrDHudUZg3qX4waLG5M43q7Wgc/MbQITxOUSQv7c7ugFFDzQGBzZswY6786m86gpE\nIbb3OhjZnzcvQAaRHhdlQWIMm2nrAgMBAAEwDQYJKoZIhvcNAQEFBQADgYEAtCu4\nnUhVVxYUntneD9+h8Mg9q6q+auNKyExzyLwaxlAoo7TJHidbtS4J5iNmZgXL0Fkb\nFFBjvSfpJIlJ00zbhNYS5f6GuoEDmFJl0ZxBHjJnyp378OD8uTs7fLvjx79LjSTb\nNYiytVbZPQUQ5Yaxu2jXnimvw3rrszlaEXAMPLE=\n-----END CERTIFICATE-----"
}
`
}
func testAccSamlDirectoryOptionalAttributeConfig() string {
	return `
resource "turbot_saml_directory" "test" {
	title               = "provider-test"
	parent              = "tmod:@turbot/turbot#/"
	profile_id_template = "{{profile.email}}"
	description         = "SAML Directory Testing1 updated"
	entry_point         = "https://example.com/myapp/sso/saml"
	certificate         = "-----BEGIN CERTIFICATE-----\nMIICiTCCAfICCQD6m7oRw0uXOjANBgkqhkiG9w0BAQUFADCBiDELMAkGA1UEBhMC\nVVMxCzAJBgNVBAgTAldBMRAwDgYDVQQHEwdTZWF0dGxlMQ8wDQYDVQQKEwZBbWF6\nb24xFDASBgNVBAsTC0lBTSBDb25zb2xlMRIwEAYDVQQDEwlUZXN0Q2lsYWMxHzAd\nBgkqhkiG9w0BCQEWEG5vb25lQGFtYXpvbi5jb20wHhcNMTEwNDI1MjA0NTIxWhcN\nMTIwNDI0MjA0NTIxWjCBiDELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAldBMRAwDgYD\nVQQHEwdTZWF0dGxlMQ8wDQYDVQQKEwZBbWF6b24xFDASBgNVBAsTC0lBTSBDb25z\nb2xlMRIwEAYDVQQDEwlUZXN0Q2lsYWMxHzAdBgkqhkiG9w0BCQEWEG5vb25lQGFt\nYXpvbi5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMaK0dn+a4GmWIWJ\n21uUSfwfEvySWtC2XADZ4nB+BLYgVIk60CpiwsZ3G93vUEIO3IyNoH/f0wYK8m9T\nrDHudUZg3qX4waLG5M43q7Wgc/MbQITxOUSQv7c7ugFFDzQGBzZswY6786m86gpE\nIbb3OhjZnzcvQAaRHhdlQWIMm2nrAgMBAAEwDQYJKoZIhvcNAQEFBQADgYEAtCu4\nnUhVVxYUntneD9+h8Mg9q6q+auNKyExzyLwaxlAoo7TJHidbtS4J5iNmZgXL0Fkb\nFFBjvSfpJIlJ00zbhNYS5f6GuoEDmFJl0ZxBHjJnyp378OD8uTs7fLvjx79LjSTb\nNYiytVbZPQUQ5Yaxu2jXnimvw3rrszlaEXAMPLE=\n-----END CERTIFICATE-----"
	name_id_format		= "UNSPECIFIED"
	allow_group_syncing = true
}
`
}

// helper functions
func testAccCheckSamlDirectoryExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadSamlDirectory(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckSamlDirectoryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_saml_directory" {
			_, err := client.ReadSamlDirectory(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Alert still exists")
			}
			if !errors.NotFoundError(err) {
				return fmt.Errorf("expected 'not found' error, got %s", err)
			}
		}
	}

	return nil
}

// samlDirectoryResponse is a SAML directory as Guardrails returns it, captured live from a workspace with the
// certificate shortened. signRequests is a Boolean, signaturePrivateKey is a reference to the stored secret rather than
// the key, and signatureAlgorithm is the stored "SHA-256" form, not the "SHA256" that the create and update inputs take.
func samlDirectoryResponse(signRequests bool, signatureAlgorithm string) string {
	return fmt.Sprintf(`{"data":{"resource":{"type":{"uri":"tmod:@turbot/turbot-iam#/resource/types/samlDirectory"},`+
		`"parent":"173245961327617","title":"saml","description":null,"status":"Active","directoryType":"saml",`+
		`"profileIdTemplate":"{{profile.email}}","entryPoint":"https://idp.example.com/sso","certificate":"CERT",`+
		`"issuer":null,"nameIdFormat":"unspecified","signRequests":%t,`+
		`"signaturePrivateKey":{"secret":{"id":"401749014329411"}},"signatureAlgorithm":%q,"poolId":null,`+
		`"allowIdpInitiatedSso":null,"allowGroupSyncing":null,"profileGroupsAttribute":null,"groupFilter":null,`+
		`"turbot":{"id":"401749014378596","parentId":"173245961327617"}}}}`, signRequests, signatureAlgorithm)
}

// samlTestClient answers every SAML directory query and mutation with directory, and the parent akas lookup that
// follows each of them. It returns the variables of each request it receives.
func samlTestClient(t *testing.T, directory string) (*apiClient.Client, func() []map[string]interface{}) {
	t.Helper()
	var mu sync.Mutex
	var received []map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Query     string
			Variables map[string]interface{}
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decoding request: %s", err)
		}
		mu.Lock()
		received = append(received, request.Variables)
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(request.Query, "signRequests") {
			fmt.Fprint(w, directory)
			return
		}
		fmt.Fprint(w, `{"data":{"resource":{"turbot":{"id":"173245961327617","akas":["tmod:@turbot/turbot#/"]}}}}`)
	}))
	t.Cleanup(server.Close)
	requests := func() []map[string]interface{} {
		mu.Lock()
		defer mu.Unlock()
		return append([]map[string]interface{}(nil), received...)
	}
	return &apiClient.Client{Graphql: graphql.NewClient(server.URL + "/graphql")}, requests
}

func samlTestData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	config := map[string]interface{}{
		"parent":              "tmod:@turbot/turbot#/",
		"title":               "saml",
		"entry_point":         "https://idp.example.com/sso",
		"certificate":         "CERT",
		"profile_id_template": "{{profile.email}}",
	}
	for k, v := range raw {
		config[k] = v
	}
	return schema.TestResourceDataRaw(t, resourceTurbotSamlDirectory().Schema, config)
}

func TestSamlDirectoryReadAcceptsGuardrailsShapes(t *testing.T) {
	for _, tc := range []struct {
		name                   string
		signRequests           bool
		storedAlgorithm        string
		wantSignRequests       string
		wantSignatureAlgorithm string
	}{
		{"signing on", true, "SHA-256", "true", "SHA256"},
		// The console's edit form always saves signRequests and signatureAlgorithm, defaulting to false and SHA1.
		{"saved in the console", false, "SHA-1", "false", "SHA1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client, _ := samlTestClient(t, samlDirectoryResponse(tc.signRequests, tc.storedAlgorithm))
			d := samlTestData(t, map[string]interface{}{"signature_private_key": "KEY-FROM-CONFIG"})
			d.SetId("401749014378596")

			if err := resourceTurbotSamlDirectoryRead(d, client); err != nil {
				t.Fatalf("read: %s", err)
			}

			assertAttr(t, d, "sign_requests", tc.wantSignRequests)
			assertAttr(t, d, "signature_algorithm", tc.wantSignatureAlgorithm)
			// Guardrails returns only a reference to the key, so state keeps the configured key.
			assertAttr(t, d, "signature_private_key", "KEY-FROM-CONFIG")
		})
	}
}

func TestSamlDirectoryWritesSendSigningAsGuardrailsTypes(t *testing.T) {
	for _, tc := range []struct {
		name  string
		write func(*schema.ResourceData, interface{}) error
	}{
		{"create", resourceTurbotSamlDirectoryCreate},
		{"update", resourceTurbotSamlDirectoryUpdate},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client, requests := samlTestClient(t, samlDirectoryResponse(true, "SHA-256"))
			d := samlTestData(t, map[string]interface{}{
				"sign_requests":         "true",
				"signature_algorithm":   "SHA256",
				"signature_private_key": "KEY-FROM-CONFIG",
			})
			d.SetId("401749014378596")

			if err := tc.write(d, client); err != nil {
				t.Fatalf("%s: %s", tc.name, err)
			}

			input, _ := requests()[0]["input"].(map[string]interface{})
			if input["signRequests"] != true {
				t.Errorf("signRequests must be sent as the Boolean true, got %#v", input["signRequests"])
			}
			if input["signatureAlgorithm"] != "SHA256" {
				t.Errorf("signatureAlgorithm must be sent as SHA256, got %#v", input["signatureAlgorithm"])
			}
			if input["signaturePrivateKey"] != "KEY-FROM-CONFIG" {
				t.Errorf("signaturePrivateKey must be sent from config, got %#v", input["signaturePrivateKey"])
			}
			assertAttr(t, d, "sign_requests", "true")
			assertAttr(t, d, "signature_algorithm", "SHA256")
			assertAttr(t, d, "signature_private_key", "KEY-FROM-CONFIG")
		})
	}
}

func TestSamlDirectoryCreateOmitsUnsetSignRequests(t *testing.T) {
	client, requests := samlTestClient(t, samlDirectoryResponse(false, "SHA-1"))
	d := samlTestData(t, nil)

	if err := resourceTurbotSamlDirectoryCreate(d, client); err != nil {
		t.Fatalf("create: %s", err)
	}

	input, _ := requests()[0]["input"].(map[string]interface{})
	if _, sent := input["signRequests"]; sent {
		t.Errorf("an unset sign_requests must not be sent, got %#v", input["signRequests"])
	}
}

func TestSamlDirectorySigningSchema(t *testing.T) {
	s := resourceTurbotSamlDirectory().Schema
	for _, name := range []string{"sign_requests", "signature_algorithm"} {
		// sign_requests holds a Boolean, but it stays a string: changing its type would break every existing state,
		// which holds "", and would stop older provider versions from reading state written by this one.
		if s[name].Type != schema.TypeString {
			t.Errorf("%s must stay a string", name)
		}
		// Computed, so a value saved in the console does not show as a diff for a configuration that leaves it unset.
		if !s[name].Optional || !s[name].Computed {
			t.Errorf("%s must be Optional and Computed", name)
		}
	}
	for _, tc := range []struct {
		name   string
		values []string
		valid  bool
	}{
		{"sign_requests", []string{"true", "false"}, true},
		{"sign_requests", []string{"Enabled", "Disabled", "TRUE", ""}, false},
		{"signature_algorithm", []string{"SHA1", "SHA256", "SHA512"}, true},
		{"signature_algorithm", []string{"SHA-256", "sha256", "SHA384", ""}, false},
	} {
		if s[tc.name].ValidateFunc == nil {
			t.Fatalf("%s must have a ValidateFunc", tc.name)
		}
		for _, v := range tc.values {
			_, errs := s[tc.name].ValidateFunc(v, tc.name)
			if (len(errs) == 0) != tc.valid {
				t.Errorf("%s = %q: valid must be %t, errors %v", tc.name, v, tc.valid, errs)
			}
		}
	}
}

func assertAttr(t *testing.T, d *schema.ResourceData, name, want string) {
	t.Helper()
	if got := d.Get(name).(string); got != want {
		t.Errorf("%s = %q, want %q", name, got, want)
	}
}
