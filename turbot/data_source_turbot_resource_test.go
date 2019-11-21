package turbot

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

// todo test more policy formats: array, templated, calculated (e.g. stack source)

func TestAccResourceDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.turbot_resource.test_resource", "turbot.title", "provider_test"),
				),
			},
		},
	})

}
func testAccResourceDataSourceConfig() string {
	return `
resource "turbot_folder" "test" {
	parent = "tmod:@turbot/turbot#/"
	title = "provider_test"
	description = "test folder for turbot terraform provider"
}

data "turbot_resource" "test_resource" {
  id = turbot_folder.test.id
}
`
}
