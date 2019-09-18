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
						"data.turbot_resource.test_resource", "turbot", "alpha.turbot.com"),
				),
			},
		},
	})

}
func testAccResourceDataSourceConfig() string {
	return `
data "turbot_resource" "test_resource" {
  aka = "tmod:@turbot/turbot#/"
}
`
}
