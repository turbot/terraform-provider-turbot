package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"github.com/terraform-providers/terraform-provider-turbot/errors"
	"testing"
)

// test suites
func TestAccLocalDirectoryUser_Basic(t *testing.T) {
	resourceName := "turbot_local_directory_user.test_user"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocalDirectoryUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLocalDirectoryUserConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalDirectoryUserExists("turbot_local_directory_user.test_user"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory_user.test_user", "title", "Kai Daguerre"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory_user.test_user", "email", "kai@turbot.com"),
				),
			},
			{
				Config: testAccLocalDirectoryUserUpdateEmailConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalDirectoryUserExists("turbot_local_directory_user.test_user"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory_user.test_user", "title", "Kai Daguerre"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory_user.test_user", "email", "kai2@turbot.com"),
				),
			},
			{
				Config: testAccLocalDirectoryUserUpdateTitleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalDirectoryUserExists("turbot_local_directory_user.test_user"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory_user.test_user", "title", "Kai Daguerre2"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory_user.test_user", "email", "kai@turbot.com"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

// configs
func testAccLocalDirectoryUserConfig() string {
	return `
resource "turbot_local_directory_user" "test_user" {
	title        = "Kai Daguerre"
	email        = "kai@turbot.com"
	display_name = "Kai Daguerre"
	parent       = "184298093985240"
}
`
}

func testAccLocalDirectoryUserUpdateTitleConfig() string {
	return `
resource "turbot_local_directory_user" "test_user" {
	title        = "Kai Daguerre2"
	email        = "kai@turbot.com"
	display_name = "Kai Daguerre"
	parent       = "184298093985240"
}`
}

func testAccLocalDirectoryUserUpdateEmailConfig() string {
	return `
resource "turbot_local_directory_user" "test_user" {
	title        = "Kai Daguerre"
	email        = "kai2@turbot.com"
	display_name = "Kai Daguerre"
	parent       = "184298093985240"
}`
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
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadLocalDirectoryUser(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckLocalDirectoryUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_local_directory_user" {
			_, err := client.ReadLocalDirectoryUser(rs.Primary.ID)
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
