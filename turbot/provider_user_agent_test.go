package turbot

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
)

// Configuring the provider the way Terraform does sends Terraform's version in the User-Agent.
func TestConfigureSendsTerraformVersionInUserAgent(t *testing.T) {
	saved := ProviderVersion
	defer func() { ProviderVersion = saved }()
	ProviderVersion = "1.16.0"
	t.Setenv("TF_APPEND_USER_AGENT", "")

	var seen string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{}}`))
	}))
	defer server.Close()
	// The workspace URL must be https, so trust the test server's certificate.
	savedTransport := http.DefaultTransport
	http.DefaultTransport = server.Client().Transport
	defer func() { http.DefaultTransport = savedTransport }()

	p := Provider().(*schema.Provider)
	p.TerraformVersion = "1.9.0"
	d := schema.TestResourceDataRaw(t, p.Schema, map[string]interface{}{
		"access_key": "a",
		"secret_key": "s",
		"workspace":  server.URL + "/api/latest",
	})
	p.ConfigureFunc(d) // Validate fails on the empty answer; only the header matters here.

	assert.Equal(t,
		"HashiCorp Terraform/1.9.0 (+https://www.terraform.io) terraform-provider-turbot/1.16.0",
		seen)
}

func TestUserAgentNamesTerraformAndProvider(t *testing.T) {
	saved := ProviderVersion
	defer func() { ProviderVersion = saved }()
	ProviderVersion = "1.16.0"
	t.Setenv("TF_APPEND_USER_AGENT", "")

	assert.Equal(t,
		"HashiCorp Terraform/1.9.0 (+https://www.terraform.io) terraform-provider-turbot/1.16.0",
		userAgent("1.9.0"))
}

func TestUserAgentWithoutTerraformVersion(t *testing.T) {
	saved := ProviderVersion
	defer func() { ProviderVersion = saved }()
	ProviderVersion = "1.16.0"
	t.Setenv("TF_APPEND_USER_AGENT", "")

	assert.Equal(t,
		"HashiCorp Terraform/unknown (+https://www.terraform.io) terraform-provider-turbot/1.16.0",
		userAgent(""))
}
