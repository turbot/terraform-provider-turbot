package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"strings"
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
				Config: testAccPolicySettingStringConfig(stringPolicyType, "testValue", "must"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "testValue"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "must"),
				),
			},
			{
				Config: testAccPolicySettingStringConfig(stringPolicyType, "testValue-updated", "must"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "testValue-updated"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "must"),
				),
			}, {
				Config: testAccPolicySettingStringConfig(stringPolicyType, "testValue-updated", "should"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "testValue-updated"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "should"),
				),
			},

			{
				Config: testAccPolicySettingTemplateConfig(stringPolicyType, stringPolicyTemplate, stringPolicyTemplateInput, "should"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "template", stringPolicyTemplate),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "template_input", stringPolicyTemplateInput),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "should"),
				),
			},
		},
	})
}

func TestAccPolicySetting_Int(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicySettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicySettingIntConfig(intPolicyType, 1, "must"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "1"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "must"),
				),
			},
			{
				Config: testAccPolicySettingIntConfig(intPolicyType, 2, "must"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "2"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "must"),
				),
			},
			// value as string
			{
				Config: testAccPolicySettingStringConfig(intPolicyType, "3", "must"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "3"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "must"),
				),
			},
		},
	})
}

func TestAccPolicySetting_Array(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicySettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicySettingStringConfig(stringArrayPolicyType, "<<EOF\n- a\n- b\n- c\nEOF", "must"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", fmt.Sprintf("%v", []string{"a", "b", "c"})),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value_source", "- a\n- b\n- c\n"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "must"),
				),
			}, {
				Config: testAccPolicySettingStringConfig(stringArrayPolicyType, "<<EOF\n- b\n- a\n- d\nEOF", "must"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", fmt.Sprintf("%v", []string{"b", "a", "d"})),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value_source", "- b\n- a\n- d\n"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "must"),
				),
			},
		},
	})
}

func TestAccPolicySetting_Secret(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicySettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicySettingStringConfig(secretPolicyType, "test1", "must"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "test1"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "must"),
				),
			},
		},
	})
}

// configs
var stringPolicyType = "tmod:@turbot/provider-test#/policy/types/stringPolicy"
var intPolicyType = "tmod:@turbot/provider-test#/policy/types/integerPolicy"
var stringArrayPolicyType = "tmod:@turbot/provider-test#/policy/types/stringArrayPolicy"
var secretPolicyType = "tmod:@turbot/provider-test#/policy/types/secretPolicy"
var stringPolicyTemplate = "{% if $.account.Id == '650022101893' %}Skip{% else %}'Check: Configured'{% endif %}"
var stringPolicyTemplateInput = "{ account{ Id } }"

func testAccPolicySettingStringConfig(policyType, value string, precedence string) string {
	// wrap value in quotes if necessary (this is so we can support heredoc)
	if !strings.HasPrefix(value, `<<EOF`) {
		value = fmt.Sprintf(`"%s"`, value)
	}

	return buildConfig(policyType, value, precedence)
}

func testAccPolicySettingIntConfig(policyType string, value int, precedence string) string {
	return buildConfig(policyType, fmt.Sprintf("%d", value), precedence)
}

func testAccPolicySettingTemplateConfig(policyType, template, templateInput, precedence string) string {
	return fmt.Sprintf(`
resource "turbot_folder" "parent" {
	parent = "tmod:@turbot/turbot#/"
	title = "terraform-provider-turbot"
	description = "terraform-provider-turbot"
}
resource "turbot_policy_setting" "test_policy" {
	resource = turbot_folder.parent.id
	policy_type = "%s"
	template = "%s"
	template_input = "%s"
	precedence = "%s"
}`, policyType, template, templateInput, precedence)
}

func buildConfig(policyType, value string, precedence string) string {

	config := fmt.Sprintf(`
resource "turbot_folder" "parent" {
	parent = "tmod:@turbot/turbot#/"
	title = "provider_acceptance_tests"
	description = "Acceptance testing folder"
}
resource "turbot_policy_setting" "test_policy" {
	resource = turbot_folder.parent.id
	policy_type = "%s"
	value = %s
	precedence = "%s"
}`, policyType, value, precedence)
	return config
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
