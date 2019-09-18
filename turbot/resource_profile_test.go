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
					resource.TestCheckResourceAttr(
						"turbot_profile.test", "title", "provider_test"),
					resource.TestCheckResourceAttr(
						"turbot_profile.test", "description", "test profile"),
				),
			},
			{
				Config: testAccProfileUpdateDescConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfileExists("turbot_profile.test"),
					resource.TestCheckResourceAttr(
						"turbot_profile.test", "title", "provider_test"),
					resource.TestCheckResourceAttr(
						"turbot_profile.test", "description", "test profile for turbot terraform provider"),
				),
			},
			{
				Config: testAccProfileUpdateTitleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfileExists("turbot_profile.test"),
					resource.TestCheckResourceAttr(
						"turbot_profile.test", "title", "provider_test_upd"),
					resource.TestCheckResourceAttr(
						"turbot_profile.test", "description", "test profile for turbot terraform provider"),
				),
			},
		},
	})
}

func TestAccProfileWithDependencies(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileWithDependenciesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfileExists("turbot_profile.parent"),
					testAccCheckProfileExists("turbot_profile.child"),
					resource.TestCheckResourceAttr(
						"turbot_profile.parent", "title", "provider_test_parent"),
					resource.TestCheckResourceAttr(
						"turbot_profile.parent", "description", "parent"),
					resource.TestCheckResourceAttr(
						"turbot_profile.child", "title", "provider_test_child"),
					resource.TestCheckResourceAttr(
						"turbot_profile.child", "description", "child"),
				),
			},
			{
				Config: testAccProfileWithDependenciesUpdateDescConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfileExists("turbot_profile.parent"),
					testAccCheckProfileExists("turbot_profile.child"),
					resource.TestCheckResourceAttr(
						"turbot_profile.parent", "title", "provider_test_parent"),
					resource.TestCheckResourceAttr(
						"turbot_profile.parent", "description", "PARENT"),
					resource.TestCheckResourceAttr(
						"turbot_profile.child", "title", "provider_test_child"),
					resource.TestCheckResourceAttr(
						"turbot_profile.child", "description", "CHILD"),
				),
			},
			{
				Config: testAccProfileWithDependenciesUpdateTitleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfileExists("turbot_profile.parent"),
					testAccCheckProfileExists("turbot_profile.child"),
					resource.TestCheckResourceAttr(
						"turbot_profile.parent", "title", "PROVIDER_TEST_PARENT"),
					resource.TestCheckResourceAttr(
						"turbot_profile.parent", "description", "PARENT"),
					resource.TestCheckResourceAttr(
						"turbot_profile.child", "title", "PROVIDER_TEST_CHILD"),
					resource.TestCheckResourceAttr(
						"turbot_profile.child", "description", "CHILD"),
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

func testAccProfileUpdateDescConfig() string {
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

func testAccProfileWithDependenciesConfig() string {
	return `
resource "turbot_profile" "parent" {
  parent = "tmod:@turbot/turbot#/"
	title = "Snape"
	display_name = "Severus Snape"
	email = "severus.slytherin@hogwards.com"
	given_name = "Severus Snape"
	family_name = "Snape"
	directory_pool_id = "snapeseverus
	status = "Active"
	profile_id = "170759063660234"
}
resource "turbot_profile" "child" {
  parent = "tmod:@turbot/turbot#/"
	title = "Snape"
	display_name = "Severus Snape"
	email = "severus.slytherin@hogwards.com"
	given_name = "Severus Snape"
	family_name = "Snape"
	directory_pool_id = "snapeseverus"
	status = "Active"
	profile_id = "170759063660234"
}`
}

func testAccProfileWithDependenciesUpdateDescConfig() string {
	return `
resource "turbot_profile" "parent" {
  parent = "tmod:@turbot/turbot#/"
  title = "provider_test_parent"
  description = "PARENT"
}
resource "turbot_profile" "child" {
  parent = "${turbot_profile.parent.id}"
  title = "provider_test_child"
  description = "CHILD"
}
`
}

func testAccProfileWithDependenciesUpdateTitleConfig() string {
	return `
resource "turbot_profile" "parent" {
  parent = "tmod:@turbot/turbot#/"
  title = "PROVIDER_TEST_PARENT"
  description = "PARENT"
}
resource "turbot_profile" "child" {
  parent = "${turbot_profile.parent.id}"
  title = "PROVIDER_TEST_CHILD"
  description = "CHILD"
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
