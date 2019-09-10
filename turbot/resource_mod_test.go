package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"testing"
)

// test suites
func TestAccMod(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckModDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckModConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModExists("turbot_mod.test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "mod", "structure-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version", "5.0.0-beta.3"),
				),
			},
			{
				Config: testAccCheckModUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckModExists("turbot_mod.test"),
					testAccCheckModExists("turbot_mod.test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "mod", "structure-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version", "5.0.0-beta.4"),
				),
			},
		},
	})
}

// configs
func testAccCheckModConfig() string {
	return `
resource "turbot_mod" "test" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "structure-test"
	version = "5.0.0"
}
`
}

func testAccCheckModUpdateConfig() string {
	return `
resource "turbot_mod" "test" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "structure-test"
	version = ">=5.0.0"
}
`
}

// helper functions
func testAccCheckModExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiclient.Client)
		_, err := client.ReadMod(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckModDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mod" {
			continue
		}
		_, err := client.ReadMod(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		if !apiclient.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}

	return nil
}
