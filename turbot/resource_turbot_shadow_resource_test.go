package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"testing"
)

// test suites
func TestAccShadowResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccShadowResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckShadowResourceExists("turbot_shadow_resource.shadow_resource"),
					resource.TestCheckResourceAttr("turbot_shadow_resource.shadow_resource", "resource", "arn:aws:logs:us-east-2:713469427990:log-group:provider-test-hashicorp"),
				),
			},
		},
	})
}

// configs
func testAccShadowResourceConfig() string {
	return fmt.Sprintf(`
resource "turbot_policy_setting" "region_stack_source" {
  resource = "arn:aws::us-east-2:713469427990"
 type = "tmod:@turbot/aws#/policy/types/regionStackSource"
  value = <<EOF
resource "aws_cloudwatch_log_group" "provider" {
  name = "provider-test-hashicorp"
  tags = {
    Environment = "production"
    Application = "serviceA"
  }
}
EOF
  precedence = "REQUIRED"
}
resource "turbot_shadow_resource" "shadow_resource" {
  resource    = "arn:aws:logs:us-east-2:713469427990:log-group:provider-test-hashicorp"
  timeouts    {
    create = "5m"
}
}`)
}

// helper functions
func testAccCheckShadowResourceExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ResourceExists(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}
