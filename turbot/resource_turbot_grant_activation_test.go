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
func TestAccGrantActivate_Basic(t *testing.T) {
	resourceName := "turbot_grant_activation.test_activation"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testAccCheckLocalGrantDestroy, testAccCheckActiveGrantDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccGrantActivateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalGrantExists("turbot_grant.test_grant"),
					testAccCheckActiveGrantExists("turbot_grant_activation.test_activation"),
					resource.TestCheckResourceAttr(
						"turbot_grant.test_grant", "resource", "tmod:@turbot/turbot#/"),
					resource.TestCheckResourceAttr(
						"turbot_grant.test_grant", "type", "tmod:@turbot/turbot-iam#/permission/types/turbot"),
					resource.TestCheckResourceAttr(
						"turbot_grant.test_grant", "level", "tmod:@turbot/turbot-iam#/permission/levels/owner"),
					resource.TestCheckResourceAttrSet(
						"turbot_grant_activation.test_activation", "resource"),
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
func testAccGrantActivateConfig() string {
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

resource "turbot_grant_activation" "test_activation" {
	resource = turbot_grant.test_grant.resource
	grant = turbot_grant.test_grant.id
}
`
}

// helper functions
func testAccCheckLocalGrantExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadGrant(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckLocalGrantDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_grant" {
			_, err := client.ReadGrant(rs.Primary.ID)
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

func testAccCheckActiveGrantExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadGrantActivation(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckActiveGrantDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		// rs.Type is the TERRAFORM type name, never the Guardrails "Grant" - the previous guard
		// never matched, so the activation's destruction was never actually verified. This is the
		// grant ACTIVATION destroy, so read the activation, not the grant.
		if rs.Type != "turbot_grant_activation" {
			continue
		}
		_, err := client.ReadGrantActivation(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("grant activation %s still exists after destroy", rs.Primary.ID)
		}
		if !errors.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}

	return nil
}
