package turbot

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccGrant_Basic(t *testing.T) {
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
						"turbot_grant.test_grant", "type", "tmod:@turbot/aws#/permission/types/aws"),
					resource.TestCheckResourceAttr(
						"turbot_grant.test_grant", "level", "tmod:@turbot/turbot-iam#/permission/levels/superuser"),
				),
			},
		},
	})
}

// configs
func testAccGrantConfig() string {
	return `
resource "turbot_local_directory" "test_dir" {
	parent              = "tmod:@turbot/turbot#/"
	title               = "provider_test_directory"
	description         = "provider_test_directory"
	profile_id_template = "{{profile.email}}"
}

resource "turbot_local_directory_user" "test_user" {
	title        = "Rupesh"
	email        = "rupesh@turbot.com"
	display_name = "rupesh"
	parent       = turbot_local_directory.test_dir.id
}

resource "turbot_profile" "test_profile" {
	title             = turbot_local_directory_user.test_user.title
	email             = turbot_local_directory_user.test_user.email
	directory_pool_id = "dpi"
	given_name 		  = "rupesh"
	family_name       = "patil"
	display_name      = turbot_local_directory_user.test_user.display_name
	parent            = turbot_local_directory.test_dir.id
	profile_id        = turbot_local_directory_user.test_user.email
}

resource "turbot_grant" "test_grant" {
	resource         = "tmod:@turbot/turbot#/"
	type  = "tmod:@turbot/turbot-iam#/permission/types/turbot"
	level = "tmod:@turbot/turbot-iam#/permission/levels/owner"
	identity          = turbot_profile.test_profile.id
}
`
}
