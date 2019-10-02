package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"testing"
)

// test suites
func TestAccResourceFolder(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfig(folderType, folderBody),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.test"),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "type", folderType),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "body", folderBody),
				),
			},
			// TODO this fails as when upserting an existing folder a new folder is created
			{
				Config: testAccResourceConfig(folderType, folderBodyUpdatedDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.test"),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "type", folderType),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "body", folderBodyUpdatedDescription),
				),
			},
			{
				Config: testAccResourceConfig(folderType, folderBodyUpdatedTitle),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("turbot_resource.test"),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "type", folderType),
					resource.TestCheckResourceAttr(
						"turbot_resource.test", "body", folderBodyUpdatedTitle),
				),
			},
		},
	})
}

// configs
var folderType = `tmod:@turbot/turbot#/resource/types/folder`
var folderBody = `{
 "title": "provider_test",
 "description": "test resource"
}
`
var folderBodyUpdatedTitle = `{
 "title": "provider_test_",
 "description": "test resource"
}
`
var folderBodyUpdatedDescription = `{
 "title": "provider_test_",
 "description": "test resource"
}
`

func testAccResourceConfig(resourceType, body string) string {
	return fmt.Sprintf(`
resource "turbot_resource" "test" {
 parent = "tmod:@turbot/turbot#/"
 type = "%s"
 body =  <<EOF
%sEOF
}
`, resourceType, body)
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
		client := testAccProvider.Meta().(*apiclient.Client)
		_, err := client.ReadResource(rs.Primary.ID, nil)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckResourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "resource" {
			continue
		}
		_, err := client.ReadResource(rs.Primary.ID, nil)
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		if !apiclient.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}

	return nil
}
