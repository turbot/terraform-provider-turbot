package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
	"regexp"
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
  filter = "resourceType:181381985925765 $.turbot.tags.a:b"
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
		// Verify the ATTACHMENT from the resource side, matching what Exists() does. Reading the
		// smart folder only proves the folder exists, and needs a grant on the folder itself.
		attached, err := client.PolicyPackAttached(resource, smartFolderId)
		if err != nil {
			return fmt.Errorf("error fetching attachment for resource %s. %s", resource, err)
		}
		if !attached {
			return fmt.Errorf("smart folder %s is not attached to resource %s", smartFolderId, resource)
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
		if !errors.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}
	return nil
}

// TestAccSmartFolderAttachment_IdAka covers the legacy twin with the smart folder given as an
// aka as well as an id. Both must attach, and smart_folder must settle to the numeric id.
func TestAccSmartFolderAttachment_IdAka(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSmartFolderAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSmartFolderAttachmentIdAkaConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartFolderAttachmentExists("turbot_smart_folder_attachment.by_id"),
					testAccCheckSmartFolderAttachmentExists("turbot_smart_folder_attachment.by_aka"),
					resource.TestMatchResourceAttr("turbot_smart_folder_attachment.by_aka", "smart_folder",
						regexp.MustCompile(`^[0-9]+$`)),
				),
			},
			{
				Config:   testAccSmartFolderAttachmentIdAkaConfig(),
				PlanOnly: true,
			},
		},
	})
}

func testAccSmartFolderAttachmentIdAkaConfig() string {
	// Each combination needs its OWN target folder: attaching the same smart folder to the
	// same resource twice is a duplicate attachment, not a second test case.
	return `
resource "turbot_smart_folder" "idaka" {
  parent      = "tmod:@turbot/turbot#/"
  filter      = "resourceType:181381985925765 $.turbot.tags.a:b"
  description = "Smart Folder id/aka Testing"
  title       = "smart_folder_idaka"
  akas        = ["test_smart_folder_idaka_aka"]
}

resource "turbot_folder" "idaka1" {
  parent      = "tmod:@turbot/turbot#/"
  title       = "provider_test_sf_idaka_1"
  description = "target for smart folder referenced by id"
}

resource "turbot_folder" "idaka2" {
  parent      = "tmod:@turbot/turbot#/"
  title       = "provider_test_sf_idaka_2"
  description = "target for smart folder referenced by aka"
}

resource "turbot_smart_folder_attachment" "by_id" {
  smart_folder = turbot_smart_folder.idaka.id
  resource     = turbot_folder.idaka1.id
}

resource "turbot_smart_folder_attachment" "by_aka" {
  smart_folder = "test_smart_folder_idaka_aka"
  resource     = turbot_folder.idaka2.id
  depends_on   = [turbot_smart_folder.idaka]
}
`
}
