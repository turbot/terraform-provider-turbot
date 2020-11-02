package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"github.com/terraform-providers/terraform-provider-turbot/errorHandler"
	"github.com/terraform-providers/terraform-provider-turbot/helpers"
	"testing"
)

// test suites
func TestAccResource_BasicFolderResource(t *testing.T) {
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
			{
				Config: testAccResourceConfigFolder(folderType, folderWithNoDescription, metadataUpdated),
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
func TestAccResource_AccountResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfigAccount(accountType, metadata, accountData),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.account_resource"),
					resource.TestCheckResourceAttr(
						"turbot_resource.account_resource", "type", accountType),
					resource.TestCheckResourceAttr(
						"turbot_resource.account_resource", "data", helpers.FormatJson(accountData)),
					resource.TestCheckResourceAttr(
						"turbot_resource.account_resource", "metadata", helpers.FormatJson(metadata)),
				),
			},
			{
				Config: testAccResourceConfigAccount(accountType, metadata, accountDataUpdatedTitle),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.account_resource"),
					resource.TestCheckResourceAttr(
						"turbot_resource.account_resource", "type", accountType),
					resource.TestCheckResourceAttr(
						"turbot_resource.account_resource", "data", helpers.FormatJson(accountDataUpdatedTitle)),
				),
			},
		},
	})
}
func TestAccResource_AccountResourceWithFullData(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfigAccount(accountType, metadata, accountData),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.account_resource"),
					resource.TestCheckResourceAttr(
						"turbot_resource.account_resource", "type", accountType),
					resource.TestCheckResourceAttr(
						"turbot_resource.account_resource", "full_data", helpers.FormatJson(fullAccountData)),
					resource.TestCheckResourceAttr(
						"turbot_resource.account_resource", "metadata", helpers.FormatJson(metadata)),
				),
			},
		},
	})
}
func TestAccResource_AccountResourceWithFullMetadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfigAccount(accountType, metadata, accountData),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.account_resource"),
					resource.TestCheckResourceAttr(
						"turbot_resource.account_resource", "type", accountType),
					resource.TestCheckResourceAttr(
						"turbot_resource.account_resource", "full_data", helpers.FormatJson(fullAccountData)),
					resource.TestCheckResourceAttr(
						"turbot_resource.account_resource", "full_metadat", helpers.FormatJson(metadata)),
				),
			},
		},
	})
}

// configs
var folderType = `tmod:@turbot/turbot#/resource/types/folder`
var accountType = `tmod:@turbot/aws#/resource/types/account`

var fullAccountData = `{
 "Id": "112233445566",
 "title": "account"
 "description": "full data account"
}
`
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
var fullAccountMetadata = `{
  "c1": "custom1",
  "c2": "custom2",
  "c3": "custom3"
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
var folderWithNoDescription = `{
 "title": "provider_test_updated"
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

func testAccResourceConfigAccount(resourceType, metadata, data string) string {
	config := fmt.Sprintf(`
resource "turbot_folder" "test" {
	parent = "tmod:@turbot/turbot#/"
	title = "account_import"
	description = "test folder for turbot terraform provider"
}
resource "turbot_resource" "account_resource" {
  parent     = turbot_folder.test.id
  type       = "%s"
  data       = <<EOF
%sEOF
  metadata = <<EOF
%sEOF
}`, resourceType, data, metadata)
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
			if !errorHandler.NotFoundError(err) {
				return fmt.Errorf("expected 'not found' error, got %s", err)
			}
		}
	}

	return nil
}
