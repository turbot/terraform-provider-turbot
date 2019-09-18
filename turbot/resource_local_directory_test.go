package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"testing"
)

// test suites
func TestAccLocalDirectory(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocalDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLocalDirectoryConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalDirectoryExists("turbot_local_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory.test", "title", "provider_test"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory.test", "description", "test Directory"),
				),
			},
			{
				Config: testAccDirectoryUpdateTitleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalDirectoryExists("turbot_local_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory.test", "title", "provider_test_refactor"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory.test", "description", "test Directory"),
				),
			},
			{
				Config: testAccLocalDirectoryUpdateDescConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalDirectoryExists("turbot_local_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory.test", "title", "provider_test"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory.test", "description", "test Directory for turbot terraform provider"),
				),
			},
		},
	})
}

// configs
func testAccLocalDirectoryConfig() string {
	return `
resource "turbot_local_directory" "test" {
	parent = "tmod:@turbot/turbot#/"
	title = "provider_test"
	description = "test Directory"
	profile_id_template = "{{profile.email}}"
}
`
}

func testAccLocalDirectoryUpdateDescConfig() string {
	return `
resource "turbot_local_directory" "test" {
	parent = "tmod:@turbot/turbot#/"
	title = "provider_test"
	description = "test Directory for turbot terraform provider"
	profile_id_template = "{{profile.email}}"
}
`
}

func testAccDirectoryUpdateTitleConfig() string {
	return `
resource "turbot_local_directory" "test" {
	parent = "tmod:@turbot/turbot#/"
	title = "provider_test_refactor"
	description = "test Directory"
	profile_id_template = "{{profile.email}}"
}
`
}

// helper functions
func testAccCheckLocalDirectoryExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiclient.Client)
		_, err := client.ReadLocalDirectory(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckLocalDirectoryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "Directory" {
			continue
		}
		_, err := client.ReadLocalDirectory(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		if !apiclient.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}

	return nil
}
