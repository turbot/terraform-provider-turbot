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
func TestAccFileResourcefile_Basic(t *testing.T) {
	resourceName := "turbot_file.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFileResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFileResourceConfigfile(fileData, filemetadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFileResourceExists(resourceName),
					resource.TestCheckResourceAttr(
						"turbot_file.test", "data", helpers.FormatJson(fileData)),
					resource.TestCheckResourceAttr(
						"turbot_file.test", "metadata", helpers.FormatJson(filemetadata)),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
			{
				Config: testAccFileResourceConfigfile(fileDataUpdatedDescription, metadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFileResourceExists(resourceName),
					resource.TestCheckResourceAttr(
						"turbot_file.test", "data", helpers.FormatJson(fileDataUpdatedDescription)),
					resource.TestCheckResourceAttr(
						"turbot_file.test", "metadata", helpers.FormatJson(metadata)),
				),
			},
			{
				Config: testAccFileResourceConfigfile(fileDataUpdatedTitle, metadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFileResourceExists(resourceName),
					resource.TestCheckResourceAttr(
						"turbot_file.test", "data", helpers.FormatJson(fileDataUpdatedTitle)),
					resource.TestCheckResourceAttr(
						"turbot_file.test", "metadata", helpers.FormatJson(metadata)),
				),
			},
			{
				Config: testAccFileResourceConfigfile(fileDataUpdatedTitle, metadataUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFileResourceExists(resourceName),
					resource.TestCheckResourceAttr(
						"turbot_file.test", "data", helpers.FormatJson(fileDataUpdatedTitle)),
					resource.TestCheckResourceAttr(
						"turbot_file.test", "metadata", helpers.FormatJson(metadataUpdated)),
				),
			},
		},
	})
}

var fileData = `{
 "title": "provider_test",
 "description": "test resource"
}
`
var filemetadata = `{
 "c1": "custom1",
 "c2": "custom2"
}
`
var filemetadataUpdated = `{
 "c1": "custom1",
 "c2": "custom3"
}
`
var fileDataUpdatedTitle = `{
 "title": "provider_test_updated",
 "description": "test resource"
}
`
var fileDataUpdatedDescription = `{
 "title": "provider_test",
 "description": "test resource_updated"
}
`

// configs
func testAccFileResourceConfigfile(data, metadata string) string {
	config := fmt.Sprintf(`
resource "turbot_file" "test" {
	parent = "tmod:@turbot/turbot#/"
	title  = "provider_file"
description = "test"
	data =  <<EOF
%sEOF
	metadata =  <<EOF
%sEOF
}
`, data, metadata)
	return config
}

// helper functions
func testAccCheckFileResourceExists(resource string) resource.TestCheckFunc {
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

func testAccCheckFileResourceDestroy(s *terraform.State) error {
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
