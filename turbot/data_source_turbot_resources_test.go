package turbot

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// todo test more policy formats: array, templated, calculated (e.g. stack source)

func TestAccResourcesDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.turbot_resources.test_resource", "filter", "tags:Name=terraform-test resourceType:folder"),
				),
			},
		},
	})
}

func testAccResourcesDataSourceConfig() string {
	return `
resource "turbot_folder" "test" {
	parent      = "tmod:@turbot/turbot#/"
	title       = "provider_test"
	description = "test folder for guardrails terraform provider"

	tags = {
		Name = "terraform-test"
	}
}

data "turbot_resources" "test_resource" {
  filter = "tags:Name=terraform-test resourceType:folder"
}
`
}
