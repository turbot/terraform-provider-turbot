package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"os"
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
	if v := os.Getenv("TURBOT_ACCESS_KEY_ID"); v == "" {
		t.Fatal("TURBOT_ACCESS_KEY_ID must be set for acceptance tests")
	}
	if v := os.Getenv("TURBOT_SECRET_ACCESS_KEY"); v == "" {
		t.Fatal("TURBOT_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("TURBOT_WORKSPACE"); v == "" {
		log.Fatal("TURBOT_WORKSPACE must be set for acceptance tests")
	}
}
