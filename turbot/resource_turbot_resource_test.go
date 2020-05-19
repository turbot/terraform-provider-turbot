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
	resourceName := "turbot_resource.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfigFolder(folderType, folderData, metadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.test"),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "type", folderType),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "data", helpers.FormatJson(folderData)),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "metadata", helpers.FormatJson(metadata)),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
			{
				Config: testAccResourceConfigFolder(folderType, folderDataUpdatedDescription, metadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.test"),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "type", folderType),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "data", helpers.FormatJson(folderDataUpdatedDescription)),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "metadata", helpers.FormatJson(metadata)),
				),
			},
			{
				Config: testAccResourceConfigFolder(folderType, folderDataUpdatedTitle, metadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.test"),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "type", folderType),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "data", helpers.FormatJson(folderDataUpdatedTitle)),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "metadata", helpers.FormatJson(metadata)),
				),
			},
			{
				Config: testAccResourceConfigFolder(folderType, folderDataUpdatedTitle, metadataUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.test"),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "type", folderType),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "data", helpers.FormatJson(folderDataUpdatedTitle)),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "metadata", helpers.FormatJson(metadataUpdated)),
				),
			},
		},
	})
}

func TestAccResourceFolder_Account(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfigAccount(accountType, "188042518944165", metadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.account_resource"),
					resource.TestCheckResourceAttr(
						"turbot_resource.account_resource", "type", accountType),
				),
			},
			{
				Config: testAccResourceConfigAccount(accountType, "192355801560817", metadata),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.account_resource"),
					resource.TestCheckResourceAttr(
						"turbot_resource.account_resource", "type", accountType),
					resource.TestCheckResourceAttr(
						"turbot_resource.account_resource", "parent", "192355801560817"),
				),
			},
		},
	})
}

// configs
var folderType = `tmod:@turbot/turbot#/resource/types/folder`
var accountType = `tmod:@turbot/aws#/resource/types/account`

var accountData = `{
 "Id": "112233445566",
 "title": "account"
}
`

var accountDataUpdatedTitle = `{
 "Id": "112233445566",
 "title": "account_updated"
}
`

var folderData = `{
 "title": "provider_test",
 "description": "test resource"
}
`
var metadata = `{
 "c1": "custom1",
 "c2": "custom2"
}
`
var metadataUpdated = `{
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

func testAccResourceConfigFolder(resourceType, data, metadata string) string {
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
func testAccResourceConfigAccount(resourceType, parent, metadata string) string {
	config := fmt.Sprintf(`
resource "turbot_resource" "account_resource" {
  parent     = "%s"
  type       = "%s"
  data       = jsonencode({
    "Id": "786233995633"
  })
  metadata = <<EOF
%sEOF
}`, parent, resourceType, metadata)
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
