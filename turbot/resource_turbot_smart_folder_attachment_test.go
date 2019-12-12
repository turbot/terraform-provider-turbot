package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"testing"
)

// test suites
func TestAccSmartFolderAttachment_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSmartFolderAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSmartFolderAttachmentConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartFolderAttachmentExists("turbot_smart_folder_attachment.test"),
				),
			},
		},
	})
}

// configs
func testAccSmartFolderAttachmentConfig() string {
	return `
resource "turbot_folder" "test" {
  parent = "tmod:@turbot/turbot#/"
  title = "provider_test"
  description = "test folder"
}

resource "turbot_smart_folder" "test" {
  parent  = "tmod:@turbot/turbot#/"
  filter = "resourceType:166872393063899 $.turbot.tags.a:b"
  description = "Smart Folder Testing"
  title = "smart_folder"
}

resource "turbot_smart_folder_attachment" "test" {
  resource = turbot_folder.test.id
  smart_folder = turbot_smart_folder.test.id
}
`
}

// helper functions
func testAccCheckSmartFolderAttachmentExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		smartFolderId, resource := parseSmartFolderId(rs.Primary.ID)
		_, err := client.ReadSmartFolder(smartFolderId)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}
func testAccCheckSmartFolderAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "smartFolder" {
			continue
		}
		_, err := client.ReadSmartFolder(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("alert still exists")
		}
		if !apiClient.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}
	return nil
}
