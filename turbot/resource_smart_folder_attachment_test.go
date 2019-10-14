package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"testing"
)

// test suites
func TestAccSmartFolderAttachment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSmartFolderAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSmartFolderAttachmentConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartFolderAttachmentExists("turbot_smart_folder_attachment.test"),
					resource.TestCheckResourceAttr("turbot_smart_folder_attachment.test", "resource", "167225763707951"),
					resource.TestCheckResourceAttr("turbot_smart_folder_attachment.test", "smart_folder", "171223595183451"),
				),
			},
		},
	})
}

// configs
func testAccSmartFolderAttachmentConfig() string {
	return `
    resource "turbot_smart_folder_attachment" "test" {
        resource = "167225763707951"
        smart_folder = "171222424857954"
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
		_, err := client.ReadSmartFolder(rs.Primary.ID)
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
			return fmt.Errorf("Alert still exists")
		}
		if !apiClient.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}
	return nil
}
