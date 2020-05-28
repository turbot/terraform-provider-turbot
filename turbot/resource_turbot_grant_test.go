package turbot

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccGrant_Basic(t *testing.T) {
	resourceName := "turbot_grant.test_grant"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testAccCheckLocalGrantDestroy, testAccCheckActiveGrantDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccGrantConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalGrantExists("turbot_grant.test_grant"),
					resource.TestCheckResourceAttr(
						"turbot_grant.test_grant", "resource", "tmod:@turbot/turbot#/"),
					resource.TestCheckResourceAttr(
						"turbot_grant.test_grant", "type", "tmod:@turbot/turbot-iam#/permission/types/turbot"),
					resource.TestCheckResourceAttr(
						"turbot_grant.test_grant", "level", "tmod:@turbot/turbot-iam#/permission/levels/owner"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// configs
func testAccGrantConfig() string {
	return `
resource "turbot_profile" "test_profile" {
	title             = "provider_test"
	email             = "rupesh@turbot.com"
	directory_pool_id = "dpi"
	given_name 		  = "rupesh"
	family_name       = "patil"
	display_name      = "rupesh"
	parent            = "184227597889872"
	profile_id        = "170759063660234"
}

resource "turbot_grant" "test_grant" {
	resource         = "tmod:@turbot/turbot#/"
	type  = "tmod:@turbot/turbot-iam#/permission/types/turbot"
	level = "tmod:@turbot/turbot-iam#/permission/levels/owner"
	identity          = turbot_profile.test_profile.id
}
`
}
