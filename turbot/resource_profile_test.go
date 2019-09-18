package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"testing"
)

// test suites
func TestAccProfile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfileExists("turbot_profile.test"),
					resource.TestCheckResourceAttr("turbot_profile.test", "title", "Snape"),
					resource.TestCheckResourceAttr("turbot_profile.test", "status", "Active"),
					resource.TestCheckResourceAttr("turbot_profile.test", "display_name", "Severus Snape"),
					resource.TestCheckResourceAttr("turbot_profile.test", "given_name", "Severus Snape"),
					resource.TestCheckResourceAttr("turbot_profile.test", "email", "severus.slytherin@hogwards.com"),
					resource.TestCheckResourceAttr("turbot_profile.test", "family_name", "Snape"),
					resource.TestCheckResourceAttr("turbot_profile.test", "profile_id", "170759063660234"),
					resource.TestCheckResourceAttr("turbot_profile.test", "directory_pool_id", "snapeseverus"),
				),
			},
			{
				Config: testAccProfileUpdateDispNameConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfileExists("turbot_profile.test"),
					resource.TestCheckResourceAttr("turbot_profile.test", "title", "Snape"),
					resource.TestCheckResourceAttr("turbot_profile.test", "status", "Active"),
					resource.TestCheckResourceAttr("turbot_profile.test", "display_name", "Severus M Snape"),
					resource.TestCheckResourceAttr("turbot_profile.test", "given_name", "Severus Snape"),
					resource.TestCheckResourceAttr("turbot_profile.test", "email", "severus.slytherin@hogwards.com"),
					resource.TestCheckResourceAttr("turbot_profile.test", "family_name", "Snape"),
					resource.TestCheckResourceAttr("turbot_profile.test", "profile_id", "170759063660234"),
					resource.TestCheckResourceAttr("turbot_profile.test", "directory_pool_id", "snapeseverus"),
				),
			},
		},
	})
}

// configs
func testAccProfileConfig() string {
	return `
	resource "turbot_profile" "test" {
		parent = "tmod:@turbot/turbot#/"
		title = "Snape"
		display_name = "Severus Snape"
		email = "severus.slytherin@hogwards.com"
		given_name = "Severus Snape"
		family_name = "Snape"
		directory_pool_id = "snapeseverus"
		status = "Active"
		profile_id = "170759063660234"
	}
`
}

func testAccProfileUpdateDispNameConfig() string {
	return `
	resource "turbot_profile" "test" {
		parent = "tmod:@turbot/turbot#/"
		title = "Snape"
		display_name = "Severus M Snape"
		email = "severus.slytherin@hogwards.com"
		given_name = "Severus Snape"
		family_name = "Snape"
		directory_pool_id = "snapeseverus"
		status = "Active"
		profile_id = "170759063660234"
	}
`
}

func testAccProfileUpdateTitleConfig() string {
	return `
	resource "turbot_profile" "test" {
		parent = "tmod:@turbot/turbot#/"
		title = "Snape
		display_name = "Severus Snape"
		email = "severus.slytherin@hogwards.com"
		given_name = "Severus Snape"
		family_name = "Snape"
		directory_pool_id = "snapeseverus"
		status = "Active"
		profile_id = "170759063660234"
	}
`
}

// helper functions
func testAccCheckProfileExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiclient.Client)
		_, err := client.ReadProfile(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "profile" {
			continue
		}
		_, err := client.ReadProfile(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		if !apiclient.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}

	return nil
}
