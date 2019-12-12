package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"github.com/terraform-providers/terraform-provider-turbot/helpers"
	"testing"
)

// test suites
func TestAccResourceFolder_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfig(folderType, folderData, folderMetadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.test"),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "type", folderType),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "data", helpers.FormatJson(folderData)),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "metadata", helpers.FormatJson(folderMetadata)),
				),
			},
			{
				Config: testAccResourceConfig(folderType, folderDataUpdatedDescription, folderMetadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.test"),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "type", folderType),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "data", helpers.FormatJson(folderDataUpdatedDescription)),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "metadata", helpers.FormatJson(folderMetadata)),
				),
			},
			{
				Config: testAccResourceConfig(folderType, folderDataUpdatedTitle, folderMetadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.test"),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "type", folderType),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "data", helpers.FormatJson(folderDataUpdatedTitle)),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "metadata", helpers.FormatJson(folderMetadata)),
				),
			},
			{
				Config: testAccResourceConfig(folderType, folderDataUpdatedTitle, folderMetadataUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.test"),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "type", folderType),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "data", helpers.FormatJson(folderDataUpdatedTitle)),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "metadata", helpers.FormatJson(folderMetadataUpdated)),
				),
			},
		},
	})
}

// configs
var folderType = `tmod:@turbot/turbot#/resource/types/folder`
var folderData = `{
 "title": "provider_test",
 "description": "test resource"
}
`
var folderMetadata = `{
 "c1": "custom1",
 "c2": "custom2"
}
`
var folderMetadataUpdated = `{
 "c1": "custom1",
 "c2": "custom3"
}
`
var folderDataUpdatedTitle = `{
 "title": "provider_test_updated",
 "description": "test resource"
}
`
var folderDataUpdatedDescription = `{
 "title": "provider_test",
 "description": "test resource_updated"
}
`

func testAccResourceConfig(resourceType, data, metadata string) string {
	config := fmt.Sprintf(`
resource "turbot_resource" "test" {
	parent = "tmod:@turbot/turbot#/"
	type = "%s"
	data =  <<EOF
%sEOF
	metadata =  <<EOF
%sEOF
}
`, resourceType, data, metadata)
	return config
}

// helper functions
func testAccCheckResourceExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadResource(rs.Primary.ID, nil)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckResourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_resource" {
			_, err := client.ReadResource(rs.Primary.ID, nil)
			if err == nil {
				return fmt.Errorf("Alert still exists")
			}
			if !apiClient.NotFoundError(err) {
				return fmt.Errorf("expected 'not found' error, got %s", err)
			}
		}
	}

	return nil
}
