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
func TestAccGuardrail_Basic(t *testing.T) {
	resourceName := "turbot_guardrail.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGuardrailDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGuardrailConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGuardrailExists("turbot_guardrail.test"),
					resource.TestCheckResourceAttr("turbot_guardrail.test", "title", "terraform_guardrail_test"),
					resource.TestCheckResourceAttr("turbot_guardrail.test", "description", "Guardrail Testing"),
				),
			},
			{
				Config: testAccGuardrailUpdateDescConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGuardrailExists("turbot_guardrail.test"),
					resource.TestCheckResourceAttr("turbot_guardrail.test", "title", "terraform_guardrail_test"),
					resource.TestCheckResourceAttr("turbot_guardrail.test", "description", "Guardrail Testing updated"),
					resource.TestCheckResourceAttr("turbot_guardrail.test", "controls.#", "1"),
					resource.TestCheckResourceAttr("turbot_guardrail.test", "controls.0", "tmod:@turbot/aws-s3#/control/types/encryptionInTransit"),
					resource.TestCheckResourceAttr("turbot_guardrail.test", "targets.#", "1"),
					resource.TestCheckResourceAttr("turbot_guardrail.test", "targets.0", "tmod:@turbot/aws#/resource/types/account"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

// configs
func testAccGuardrailConfig() string {
	return `
resource "turbot_guardrail" "test" {
	description = "Guardrail Testing"
	title = "terraform_guardrail_test"

	targets  = ["tmod:@turbot/aws#/resource/types/account"]
  	controls = ["tmod:@turbot/aws-s3#/control/types/encryptionInTransit"]
}
`
}

func testAccGuardrailUpdateDescConfig() string {
	return `
resource "turbot_guardrail" "test" {
	description = "Guardrail Testing updated"
	title ="terraform_guardrail_test"

	targets  = ["tmod:@turbot/aws#/resource/types/account"]
  	controls = ["tmod:@turbot/aws-s3#/control/types/encryptionInTransit"]
}
`
}

// helper functions
func testAccCheckGuardrailExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadGuardrail(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckGuardrailDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_guardrail" {
			_, err := client.ReadGuardrail(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("alert still exists")
			}
			if !errors.ForbiddenError(err) {
				return fmt.Errorf("expected 'not found' error, got %s", err)
			}
		}
	}

	return nil
}
