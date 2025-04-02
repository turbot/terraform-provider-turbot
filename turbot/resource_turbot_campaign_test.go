package turbot

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
)

// test suites
func TestAccCampaign_Basic(t *testing.T) {
	resourceName := "turbot_campaign.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCampaignDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCampaignConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCampaignExists("turbot_campaign.test"),
					resource.TestCheckResourceAttr("turbot_campaign.test", "title", "Test Campaign Resource Created Through Terraform"),
					resource.TestCheckResourceAttr("turbot_campaign.test", "description", "Campaign For Testing"),
				),
			},
			{
				Config: testAccCampaignUpdateDescConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCampaignExists("turbot_campaign.test"),
					resource.TestCheckResourceAttr("turbot_campaign.test", "title", "Test Campaign Resource Created Through Terraform"),
					resource.TestCheckResourceAttr("turbot_campaign.test", "description", "Campaign For Testing Updated"),
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
func testAccCampaignConfig() string {
	return `
resource "turbot_campaign" "test" {

  title       = "Test Campaign Resource Created Through Terraform"
  description = "Campaign For Testing"

  guardrails = ["348351678444334"]
  accounts   = ["330101957430965"]

  preview {
    start_at       = "2025-11-29T00:00:00Z"
	start_early_if = "no_alerts"
	start_notice   = "enabled"
  }
}
`
}

func testAccCampaignUpdateDescConfig() string {
	return `
resource "turbot_campaign" "test" {

  title       = "Test Campaign Resource Created Through Terraform"
  description = "Campaign For Testing Updated"

  guardrails = ["348351678444334"]
  accounts   = ["330101957430965"]

  preview {
    start_at       = "2025-11-29T00:00:00Z"
	start_early_if = "no_alerts"
	start_notice   = "enabled"
  }
}
`
}

// helper functions
func testAccCheckCampaignExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadCampaign(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckCampaignDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_campaign" {
			_, err := client.ReadCampaign(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("alert still exists")
			}
			if !errors.ForbiddenError(err) {
				return fmt.Errorf("expected 'forbidden' error, got %s", err)
			}
		}
	}

	return nil
}
