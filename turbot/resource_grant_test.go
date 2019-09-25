package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"testing"
)

// test suites
func TestAccGrant(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocalGrantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLocalGrantConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalGrantExists("turbot_grant.test"),
					resource.TestCheckResourceAttr(
						"turbot_grant.test", "resource", "165808811630593"),
					resource.TestCheckResourceAttr(
						"turbot_grant.test", "profile_id", "166523243823562"),
				),
			},
		},
	})
}

// configs
func testAccLocalGrantConfig() string {
	return `
	resource "turbot_grant" "test" {
		resource = "165808811630593"
		permission_type = "165808822449188"
		permission_level = "165808822475826"
		profile_id = "166523243823562"
	  }
`
}

// helper functions
func testAccCheckLocalGrantExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiclient.Client)
		_, err := client.ReadGrant(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckLocalGrantDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "Grant" {
			continue
		}
		_, err := client.ReadGrant(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		if !apiclient.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}

	return nil
}
