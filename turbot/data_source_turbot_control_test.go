package turbot

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccControlDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccControlConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.turbot_control.test", "type", "tmod:@turbot/turbot#/control/types/controlInstalled"),
				),
			},
		},
	})

}
func testAccControlConfig() string {
	return `
data "turbot_control" "test" {
  id = "178806515688264"
}
`
}
