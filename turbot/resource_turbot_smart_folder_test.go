package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"github.com/terraform-providers/terraform-provider-turbot/errorHandler"
	"testing"
)

// test suites
func TestAccSmartFolder_Basic(t *testing.T) {
	resourceName := "turbot_smart_folder.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSmartFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSmartFolderConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartFolderExists("turbot_smart_folder.test"),
					resource.TestCheckResourceAttr("turbot_smart_folder.test", "title", "smart_folder"),
					resource.TestCheckResourceAttr("turbot_smart_folder.test", "description", "Smart Folder Testing"),
					resource.TestCheckResourceAttr("turbot_smart_folder.test", "parent", "178806508050433"),
					resource.TestCheckResourceAttr("turbot_smart_folder.test", "filter", "resourceType:181381985925765 $.turbot.tags.a:b"),
				),
			},
			{
				Config: testAccSmartFolderUpdateDescConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartFolderExists("turbot_smart_folder.test"),
					resource.TestCheckResourceAttr("turbot_smart_folder.test", "title", "smart_folder"),
					resource.TestCheckResourceAttr("turbot_smart_folder.test", "description", "Smart Folder updated"),
					resource.TestCheckResourceAttr("turbot_smart_folder.test", "parent", "178806508050433"),
					resource.TestCheckResourceAttr("turbot_smart_folder.test", "filter", "resourceType:181381985925765 $.turbot.tags.a:b"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"filter"},
			},
		},
	})
}

// configs
func testAccSmartFolderConfig() string {
	return `
resource "turbot_smart_folder" "test" {
	parent  = "tmod:@turbot/turbot#/"
	filter = "resourceType:181381985925765 $.turbot.tags.a:b"
	description = "Smart Folder Testing"
	title = "smart_folder"
}
`
}

func testAccSmartFolderUpdateDescConfig() string {
	return `
resource "turbot_smart_folder" "test" {
	parent  = "tmod:@turbot/turbot#/"
	filter = "resourceType:181381985925765 $.turbot.tags.a:b"
	description = "Smart Folder updated"
	title ="smart_folder"
}
`
}

// helper functions
func testAccCheckSmartFolderExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadSmartFolder(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckSmartFolderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_smart_folder" {
			_, err := client.ReadSmartFolder(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("alert still exists")
			}
			if !errorHandler.NotFoundError(err) {
				return fmt.Errorf("expected 'not found' error, got %s", err)
			}
		}
	}

	return nil
}
