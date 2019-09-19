package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"testing"
)

// test suites
func TestAccPolicySetting_String(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicySettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPolicySettingStringConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "Skip"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "must"),
				),
			},
			{
				Config: testAccCheckPolicySettingStringUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "Check: Configured"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "must"),
				),
			},
			{
				Config: testAccCheckPolicySettingStringUpdatePrecedenceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "Check: Configured"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "should"),
				),
			},
			{
				Config: testAccCheckPolicySettingStringUpdateTemplateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "template", `{% if $.account.Id == '650022101893' %}Skip{% else %}'Check: Configured'{% endif %}`),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "template_input", "{ account{ Id } }"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "should"),
				),
			},
		},
	})
}

// configs
var resourceAka = "tmod:@turbot/turbot#/"
var policyTypeUri = "tmod:@turbot/aws#/policy/types/accountStack"

func testAccCheckPolicySettingStringConfig() string {
	return fmt.Sprintf(`
resource "turbot_policy_setting" "test_policy" {
  resource = "%s"
  policy_type = "%s"
  value = "Skip"
  precedence = "must"
}
`, resourceAka, policyTypeUri)
}

func testAccCheckPolicySettingStringUpdateConfig() string {
	return fmt.Sprintf(`
resource "turbot_policy_setting" "test_policy" {
  resource = "%s"
  policy_type = "%s"
  value = "Check: Configured"
  precedence = "must"
}
`, resourceAka, policyTypeUri)
}

func testAccCheckPolicySettingStringUpdatePrecedenceConfig() string {
	return fmt.Sprintf(`
resource "turbot_policy_setting" "test_policy" {
  resource = "%s"
  policy_type = "%s"
  value = "Check: Configured"
  precedence = "should"
}
`, resourceAka, policyTypeUri)
}

func testAccCheckPolicySettingStringUpdateTemplateConfig() string {
	template := `"{% if $.account.Id == '650022101893' %}Skip{% else %}'Check: Configured'{% endif %}"`
	return fmt.Sprintf(`
resource "turbot_policy_setting" "test_policy" {
  resource = "%s"
  policy_type = "%s"
  template_input = "{ account{ Id } }"
  template = %s
  precedence = "should"
}
`, resourceAka, policyTypeUri, template)
}

// helper functions
func testAccCheckPolicySettingExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiclient.Client)
		_, err := client.ReadPolicySetting(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicySettingDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "policy_setting" {
			continue
		}
		_, err := client.ReadPolicySetting(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		if !apiclient.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}

	return nil
}
