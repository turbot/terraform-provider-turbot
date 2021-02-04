package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
	"testing"
)

// test suites
func TestAccLocalDirectory_Basic(t *testing.T) {
	resourceName := "turbot_local_directory.test"
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
						"turbot_local_directory.test", "description", "test directory"),
				),
			},
			{
				Config: testAccDirectoryUpdateTitleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalDirectoryExists("turbot_local_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory.test", "title", "provider_test_refactor"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory.test", "description", "test directory"),
				),
			},
			{
				Config: testAccLocalDirectoryUpdateDescConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalDirectoryExists("turbot_local_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory.test", "title", "provider_test"),
					resource.TestCheckResourceAttr(
						"turbot_local_directory.test", "description", "test directory for turbot terraform provider"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLocalDirectory_tags(t *testing.T) {
	resourceName := "turbot_local_directory.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocalDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDirectoryTagsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalDirectoryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.tag1", "tag1value"),
					resource.TestCheckResourceAttr(resourceName, "tags.tag2", "tag2value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
	description = "test directory"
	profile_id_template = "{{profile.email}}"
}
`
}

func testAccLocalDirectoryUpdateDescConfig() string {
	return `
resource "turbot_local_directory" "test" {
	parent = "tmod:@turbot/turbot#/"
	title = "provider_test"
	description = "test directory for turbot terraform provider"
	profile_id_template = "{{profile.email}}"
}
`
}

func testAccDirectoryUpdateTitleConfig() string {
	return `
resource "turbot_local_directory" "test" {
	parent = "tmod:@turbot/turbot#/"
	title = "provider_test_refactor"
	description = "test directory"
	profile_id_template = "{{profile.email}}"
}
`
}

func testAccDirectoryTagsConfig() string {
	return `
resource "turbot_local_directory" "test" {
	parent = "tmod:@turbot/turbot#/"
	title = "provider_test_refactor"
	description = "test directory"
	profile_id_template = "{{profile.email}}"
	tags = {
		  tag1 = "tag1value"
		  tag2 = "tag2value"
	}
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
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadLocalDirectory(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckLocalDirectoryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_local_directory" {
			_, err := client.ReadLocalDirectory(rs.Primary.ID)
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
