package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"testing"
)

// test suites
func TestAccLocalDirectoryUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocalDirectoryUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLocalDirectoryUserConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalDirectoryUserExists("turbot_local_directory_user.test"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory_user.test", "title", "provider_test"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory_user.test", "email", "test@turbot.com"),
				),
			},
			{
				Config: testAccLocalDirectoryUserUpdateEmailConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalDirectoryUserExists("turbot_local_directory_user.test"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory_user.test", "title", "provider_test"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory_user.test", "email", "test2@turbot.com"),
				),
			},
			{
				Config: testAccLocalDirectoryUserUpdateTitleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalDirectoryUserExists("turbot_local_directory_user.test"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory_user.test", "title", "provider_test_upd"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory_user.test", "email", "test@turbot.com"),
				),
			},
		},
	})
}

// configs
func testAccLocalDirectoryUserConfig() string {
	return `
	resource "turbot_local_directory_user" "test" {
		title = "provider_test"
		email = "test@turbot.com"
		status = "Active"
		display_name = "ProviderTest"
		parent = "170772056456165"
} 
`
}

func testAccLocalDirectoryUserUpdateTitleConfig() string {
	return `
	resource "turbot_local_directory_user" "test" {
		title = "provider_test_upd"
		email = "test@turbot.com"
		status = "Active"
		display_name = "ProviderTest"
		parent = "170772056456165"
}
`
}

func testAccLocalDirectoryUserUpdateEmailConfig() string {
	return `
	resource "turbot_local_directory_user" "test" {
		title = "provider_test"
		email = "test2@turbot.com"
		status = "Active"
		display_name = "ProviderTest"
		parent = "170772056456165"
}
`
}

// helper functions
func testAccCheckLocalDirectoryUserExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiclient.Client)
		_, err := client.ReadLocalDirectoryUser(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckLocalDirectoryUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "localDirectoryUser" {
			continue
		}
		_, err := client.ReadLocalDirectoryUser(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		if !apiclient.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}
	return nil
}
