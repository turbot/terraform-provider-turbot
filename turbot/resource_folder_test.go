package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"testing"
)

func TestAccFolder(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFolderConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleFolderExists("turbot_folder.test"),
					resource.TestCheckResourceAttr(
						"turbot_folder.test", "title", "tf_test"),
					resource.TestCheckResourceAttr(
						"turbot_folder.test", "description", "test folder"),
				),
			},
			{
				Config: testAccCheckFolderUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleFolderExists("turbot_folder.test"),
					resource.TestCheckResourceAttr(
						"turbot_folder.test", "title", "tf_test"),
					resource.TestCheckResourceAttr(
						"turbot_folder.test", "description", "test folder for turbot terraform provider"),
				),
			},
		},
	})
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

func testAccCheckFolderConfig() string {
	return `
resource "turbot_folder" "test" {
	parent = "tmod:@turbot/turbot#/"
	title = "tf_test"
	description = "test folder"
}
`
}

func testAccCheckFolderUpdateConfig() string {
	return `
resource "turbot_folder" "test" {
	parent = "tmod:@turbot/turbot#/"
	title = "tf_test"
	description = "test folder for turbot terraform provider"
}
`
}

func testAccCheckExampleFolderExists(resource string) resource.TestCheckFunc {
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
