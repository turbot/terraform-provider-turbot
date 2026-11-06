package turbot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
