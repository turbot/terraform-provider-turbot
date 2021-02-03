package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"turbot": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	_, err := apiClient.GetCredentials(apiClient.ClientConfig{})
	if err != nil {
		t.Fatal("No credentials are set - either set TURBOT_ACCESS_KEY, TURBOT_SECRET_KEY and TURBOT_WORKSPACE or populate the file ~/.config/turbot/credentials.yml")
	}
}
