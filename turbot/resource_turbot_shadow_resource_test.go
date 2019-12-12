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
				Config: testAccShadowResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckShadowResourceExists("turbot_shadow_resource.test"),
					resource.TestCheckResourceAttr("turbot_shadow_resource.test", "filter", "resource:${aws_s3_bucket.my_bucket.arn}"),
				),
			},
		},
	})
}

// configs
func testAccShadowResourceConfig() string {
	return `
resource "aws_s3_bucket" "my_bucket" {
	bucket = "shadow-resource-test"
}
resource "turbot_shadow_resource" "shadow_resource1" {
	filter    = "resource:${aws_s3_bucket.my_bucket.arn}"
}
resource "turbot_policy_setting" "s3_bucket_versioning1" {
	resource = "${turbot_shadow_resource.shadow_resource1.id}"
	policy_type = "tmod:@turbot/aws-s3#/policy/types/bucketVersioning"
	value = "Enforce: Disabled"
	precedence = "must"
}
`
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
