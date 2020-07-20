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

func TestAccControlDataSource_TypeCheck(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccControlTypeConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.turbot_control.test", "type", "tmod:@turbot/turbot#/control/types/controlInstalled"),
				),
			},
		},
	})
}

//config
func testAccControlConfig() string {
	return `
data "turbot_control" "test" {
  id = "178806515688264"
}
`
}

func testAccControlTypeConfig() string {
	return `
data "turbot_control" "test" {
	type = "tmod:@turbot/turbot#/control/types/controlInstalled"
	resource = "178806515411691"
}
`
}
