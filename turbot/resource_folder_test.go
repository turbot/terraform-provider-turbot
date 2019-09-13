package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"testing"
)

// test suites
func TestAccFolder(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFolderConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFolderExists("turbot_folder.test"),
					resource.TestCheckResourceAttr(
						"turbot_folder.test", "title", "provider_test"),
					resource.TestCheckResourceAttr(
						"turbot_folder.test", "description", "test folder"),
				),
			},
			{
				Config: testAccFolderUpdateDescConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFolderExists("turbot_folder.test"),
					resource.TestCheckResourceAttr(
						"turbot_folder.test", "title", "provider_test"),
					resource.TestCheckResourceAttr(
						"turbot_folder.test", "description", "test folder for turbot terraform provider"),
				),
			},
			{
				Config: testAccFolderUpdateTitleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFolderExists("turbot_folder.test"),
					resource.TestCheckResourceAttr(
						"turbot_folder.test", "title", "provider_test_upd"),
					resource.TestCheckResourceAttr(
						"turbot_folder.test", "description", "test folder for turbot terraform provider"),
				),
			},
		},
	})
}

func TestAccFolderWithDependencies(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFolderWithDependenciesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFolderExists("turbot_folder.parent"),
					testAccCheckFolderExists("turbot_folder.child"),
					resource.TestCheckResourceAttr(
						"turbot_folder.parent", "title", "provider_test_parent"),
					resource.TestCheckResourceAttr(
						"turbot_folder.parent", "description", "parent"),
					resource.TestCheckResourceAttr(
						"turbot_folder.child", "title", "provider_test_child"),
					resource.TestCheckResourceAttr(
						"turbot_folder.child", "description", "child"),
				),
			},
			{
				Config: testAccFolderWithDependenciesUpdateDescConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFolderExists("turbot_folder.parent"),
					testAccCheckFolderExists("turbot_folder.child"),
					resource.TestCheckResourceAttr(
						"turbot_folder.parent", "title", "provider_test_parent"),
					resource.TestCheckResourceAttr(
						"turbot_folder.parent", "description", "PARENT"),
					resource.TestCheckResourceAttr(
						"turbot_folder.child", "title", "provider_test_child"),
					resource.TestCheckResourceAttr(
						"turbot_folder.child", "description", "CHILD"),
				),
			},
			{
				Config: testAccFolderWithDependenciesUpdateTitleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFolderExists("turbot_folder.parent"),
					testAccCheckFolderExists("turbot_folder.child"),
					resource.TestCheckResourceAttr(
						"turbot_folder.parent", "title", "PROVIDER_TEST_PARENT"),
					resource.TestCheckResourceAttr(
						"turbot_folder.parent", "description", "PARENT"),
					resource.TestCheckResourceAttr(
						"turbot_folder.child", "title", "PROVIDER_TEST_CHILD"),
					resource.TestCheckResourceAttr(
						"turbot_folder.child", "description", "CHILD"),
				),
			},
		},
	})
}

// configs
func testAccFolderConfig() string {
	return `
resource "turbot_folder" "test" {
	parent = "tmod:@turbot/turbot#/"
	title = "provider_test"
	description = "test folder"
}
`
}

func testAccFolderUpdateDescConfig() string {
	return `
resource "turbot_folder" "test" {
	parent = "tmod:@turbot/turbot#/"
	title = "provider_test"
	description = "test folder for turbot terraform provider"
}
`
}

func testAccFolderUpdateTitleConfig() string {
	return `
resource "turbot_folder" "test" {
	parent = "tmod:@turbot/turbot#/"
	title = "provider_test_upd"
	description = "test folder for turbot terraform provider"
}
`
}

func testAccFolderWithDependenciesConfig() string {
	return `
resource "turbot_folder" "parent" {
  parent = "tmod:@turbot/turbot#/"
  title = "provider_test_parent"
  description = "parent"
}
resource "turbot_folder" "child" {
  parent = "${turbot_folder.parent.id}"
  title = "provider_test_child"
  description = "child"
}`
}

func testAccFolderWithDependenciesUpdateDescConfig() string {
	return `
resource "turbot_folder" "parent" {
  parent = "tmod:@turbot/turbot#/"
  title = "provider_test_parent"
  description = "PARENT"
}
resource "turbot_folder" "child" {
  parent = "${turbot_folder.parent.id}"
  title = "provider_test_child"
  description = "CHILD"
}
`
}

func testAccFolderWithDependenciesUpdateTitleConfig() string {
	return `
resource "turbot_folder" "parent" {
  parent = "tmod:@turbot/turbot#/"
  title = "PROVIDER_TEST_PARENT"
  description = "PARENT"
}
resource "turbot_folder" "child" {
  parent = "${turbot_folder.parent.id}"
  title = "PROVIDER_TEST_CHILD"
  description = "CHILD"
}
`
}

// helper functions
func testAccCheckFolderExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiclient.Client)
		_, err := client.ReadFolder(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckFolderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "folder" {
			continue
		}
		_, err := client.ReadFolder(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		if !apiclient.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}

	return nil
}
