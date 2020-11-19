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
func TestAccTurbotDirectory_Basic(t *testing.T) {
	resourceName := "turbot_turbot_directory.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTurbotDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTurbotDirectoryConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTurbotDirectoryExists("turbot_turbot_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_turbot_directory.test", "title", "provider_test"),
					resource.TestCheckResourceAttr(
						"turbot_turbot_directory.test", "description", "test directory"),
				),
			},
			{
				Config: testAccTurbotDirectoryUpdateTitleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTurbotDirectoryExists("turbot_turbot_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_turbot_directory.test", "title", "provider_test_refactor"),
					resource.TestCheckResourceAttr(
						"turbot_turbot_directory.test", "description", "test directory"),
				),
			},
			{
				Config: testAccTurbotDirectoryTagsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTurbotDirectoryExists("turbot_turbot_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_turbot_directory.test", "title", "provider_test_refactor"),
					resource.TestCheckResourceAttr(
						"turbot_turbot_directory.test", "description", "test directory"),
					resource.TestCheckResourceAttr(
						"turbot_turbot_directory.test", "tags.%", "1"),
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
func testAccTurbotDirectoryConfig() string {
	return `
	resource "turbot_turbot_directory" "test" {
	parent = "tmod:@turbot/turbot#/"
  	title = "provider_test"
  	description = "test directory"
  	profile_id_template = "{{profile.email}}"
  	server = "test"
}
`
}
func testAccTurbotDirectoryUpdateTitleConfig() string {
	return `
	resource "turbot_turbot_directory" "test" {
	parent = "tmod:@turbot/turbot#/"
  	title = "provider_test_refactor"
  	description = "test directory"
  	profile_id_template = "{{profile.email}}"
  	server = "test"
}`
}

func testAccTurbotDirectoryTagsConfig() string {
	return `
	resource "turbot_turbot_directory" "test" {
	parent = "tmod:@turbot/turbot#/"
  	title = "provider_test_refactor"
  	description = "test directory"
  	profile_id_template = "{{profile.email}}"
  	server = "test"
	tags = {
		dev = "prod"
	}
}`
}

// helper functions
func testAccCheckTurbotDirectoryExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadTurbotDirectory(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckTurbotDirectoryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_turbot_directory" {
			_, err := client.ReadTurbotDirectory(rs.Primary.ID)
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
