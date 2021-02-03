package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
	"testing"
)

// test suites
func TestAccProfile_Basic(t *testing.T) {
	resourceName := "turbot_profile.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProfileExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", "Snape"),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Severus Snape"),
					resource.TestCheckResourceAttr(resourceName, "given_name", "Severus Snape"),
					resource.TestCheckResourceAttr(resourceName, "email", "severus.slytherin@hogwards.com"),
					resource.TestCheckResourceAttr(resourceName, "family_name", "Snape"),
					resource.TestCheckResourceAttr(resourceName, "profile_id", "170759063660234"),
					resource.TestCheckResourceAttr(resourceName, "directory_pool_id", "snapeseverus"),
					resource.TestCheckResourceAttr(resourceName, "picture", "https://lh3.googleusercontent.com/a-/AOh14Gh2rSScXBQAWCauydm0ATSoeHueEfSFv5wK8SR3"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"directory_pool_id"},
			},
		},
	})
}

// configs
func testAccProfileConfig() string {
	return `
	resource "turbot_profile" "test" {
		parent = "184298093985240"
		title = "Snape"
		display_name = "Severus Snape"
		email = "severus.slytherin@hogwards.com"
        picture = "https://lh3.googleusercontent.com/a-/AOh14Gh2rSScXBQAWCauydm0ATSoeHueEfSFv5wK8SR3"
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
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadProfile(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_profile" {
			_, err := client.ReadProfile(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Alert still exists")
			}
			if !errors.NotFoundError(err) {
				return fmt.Errorf("expected 'not found' error, got %s", err)
			}
		}
	}

	return nil
}
