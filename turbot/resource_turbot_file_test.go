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
func TestAccFileResourceFile_Basic(t *testing.T) {
	resourceName := "turbot_file.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFileResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFileResourceConfigfile(fileContent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFileResourceExists(resourceName),
					resource.TestCheckResourceAttr(
						"turbot_file.test", "content", helpers.FormatJson(fileContent)),
					resource.TestCheckResourceAttr(
						"turbot_file.test", "title", "provider_file"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
			{
				Config: testAccFileResourceConfigfile(fileContentDeleteKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFileResourceExists(resourceName),
					resource.TestCheckResourceAttr(
						"turbot_file.test", "content", helpers.FormatJson(fileContentDeleteKey)),
				),
			},
			{
				Config: testAccFileResourceConfigfile(fileContentUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFileResourceExists(resourceName),
					resource.TestCheckResourceAttr(
						"turbot_file.test", "content", helpers.FormatJson(fileContentUpdated)),
				),
			},
		},
	})
}
func TestAccFileResourceFile_NoFileContent(t *testing.T) {
		resourceName := "turbot_file.test"
		resource.Test(t, resource.TestCase{
				PreCheck:     func() { testAccPreCheck(t) },
				Providers:    testAccProviders,
				CheckDestroy: testAccCheckFileResourceDestroy,
				Steps: []resource.TestStep{
						{
							Config: testAccFileResourceWithNoContent(),
							Check: resource.ComposeTestCheckFunc(
								testAccCheckFileResourceExists(resourceName),
								resource.TestCheckResourceAttr(
									"turbot_file.test", "title", "provider_file"),
							),
						},
						{
							ResourceName: resourceName,
							ImportState:  true,
						},
					},
			})
}

var fileContent = `{
 "foo": "provider_test",
 "bar": "test resource"
}
`
var fileContentUpdated = `{
 "foo": "provider_test_updated",
 "bar": "test resource"
}
`
var fileContentDeleteKey = `{
 "foo": "provider_test"
}
`

// configs
func testAccFileResourceConfigfile(Content string) string {
	config := fmt.Sprintf(`
resource "turbot_file" "test" {
	parent = "tmod:@turbot/turbot#/"
	title  = "provider_file"
    description = "test"
	content =  <<EOF
%sEOF
}
`, Content)
	return config
}
func testAccFileResourceWithNoContent() string {
		return`
resource "turbot_file" "test" {
	parent = "tmod:@turbot/turbot#/"
	title  = "provider_file"
    description = "test"
}
`}
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
		if rs.Type == "turbot_file" {
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
