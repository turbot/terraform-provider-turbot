package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"testing"
)

func TestAccMod(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckModDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckModConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleModExists("turbot_mod.logs"),
					resource.TestCheckResourceAttr(
						"turbot_mod.logs", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.logs", "mod", "aws-logs"),
					resource.TestCheckResourceAttr(
						"turbot_mod.logs", "version", "5.0.0-beta.3"),
				),
			},
			{
				Config: testAccCheckModUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleModExists("turbot_mod.backup"),
					testAccCheckExampleModExists("turbot_mod.logs"),
					resource.TestCheckResourceAttr(
						"turbot_mod.logs", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.logs", "mod", "aws-logs"),
					resource.TestCheckResourceAttr(
						"turbot_mod.logs", "version", "5.0.0-beta.4"),
				),
			},
		},
	})
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

func testAccCheckModConfig() string {
	return `
resource "turbot_mod" "logs" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "aws-logs"
	version = "5.0.0-beta.3"
}
`
}

func testAccCheckModUpdateConfig() string {
	return `
resource "turbot_mod" "logs" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "aws-logs"
	version = "5.0.0-beta.4"
}
`
}

func testAccCheckExampleModExists(resource string) resource.TestCheckFunc {
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
