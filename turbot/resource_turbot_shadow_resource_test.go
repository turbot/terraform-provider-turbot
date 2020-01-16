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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckShadowResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccShadowResourceConfig(bucket),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckShadowResourceExists("turbot_shadow_resource.shadow_resource"),
					resource.TestCheckResourceAttr("turbot_shadow_resource.shadow_resource", "resource", "arn:aws:s3:::provider-test-hashicorp"),
				),
			},
		},
	})
}

var bucket = "provider-test-hashicorp"

// configs
func testAccShadowResourceConfig(bucket string) string {
	return fmt.Sprintf(`resource "turbot_policy_setting" "region_stack" {
  resource = "arn:aws::us-east-2:650022101893"
  type = "tmod:@turbot/aws#/policy/types/regionStack"
  value = "Enforce: Configured"
  precedence = "REQUIRED"
}
resource "turbot_policy_setting" "region_stack_source" {
  resource = "arn:aws::us-east-2:650022101893"
 type = "tmod:@turbot/aws#/policy/types/regionStackSource"
  value = <<EOF
resource "aws_s3_bucket" "b" {
  bucket = %s
  acl    = "private"
  policy = <<POLICY
  {
  "Version":"2012-10-17",
  "Statement":[
    {
      "Sid":"PublicRead",
      "Effect":"Allow",
      "Principal": "*",
      "Action":["s3:GetObject"],
      "Resource":["arn:aws:s3:::testing-007/*"]
    }
  ]
}
POLICY
EOF
  precedence = "REQUIRED"
}
resource "turbot_shadow_resource" "shadow_resource" {
  resource    = "arn:aws:s3:::provider-test-hashicorp"
}`, bucket)
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

func testAccCheckShadowResourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_shadow_resource" {
			_, err := client.ResourceExists(rs.Primary.ID)
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
