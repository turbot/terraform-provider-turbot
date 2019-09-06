package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccPolicyValueDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicySettingStringConfig(),
			},
			{
				Config: testAccCheckPolicyValueConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"turbot_policy_value.test_policy", "value", "Skip"),
					resource.TestCheckResourceAttr(
						"turbot_policy_value.test_policy", "precedence", "must"),
				),
			},
		},
	})

}
func testAccCheckPolicyValueConfig() string {
	return fmt.Sprintf(`
data "turbot_policy_value" "test_policy" {
  resource = "%s"
  policy_type = "tmod:@turbot/ssl-check#/policy/types/sslCheck"
}
`, regionResourceAka)
}
