package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
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
				Config: testAccPolicySettingStringConfig(stringPolicyType, "testValue", "REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "testValue"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "REQUIRED"),
				),
			},
			{
				Config: testAccPolicySettingStringConfig(stringPolicyType, "testValue-updated", "REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "testValue-updated"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "REQUIRED"),
				),
			}, {
				Config: testAccPolicySettingStringConfig(stringPolicyType, "testValue-updated", "RECOMMENDED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "testValue-updated"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "RECOMMENDED"),
				),
			},

			{
				Config: testAccPolicySettingTemplateConfig(stringPolicyType, stringPolicyTemplate, stringPolicyTemplateInput, "RECOMMENDED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "template", stringPolicyTemplate),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "template_input", stringPolicyTemplateInput),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "RECOMMENDED"),
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
				Config: testAccPolicySettingIntConfig(intPolicyType, 1, "REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "1"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "REQUIRED"),
				),
			},
			{
				Config: testAccPolicySettingIntConfig(intPolicyType, 2, "REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "2"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "REQUIRED"),
				),
			},
			// value as string
			{
				Config: testAccPolicySettingStringConfig(intPolicyType, "3", "REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "3"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "REQUIRED"),
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
				Config: testAccPolicySettingStringConfig(stringArrayPolicyType, "<<EOF\n- a\n- b\n- c\nEOF", "REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", fmt.Sprintf("%v", []string{"a", "b", "c"})),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value_source", "- a\n- b\n- c\n"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "REQUIRED"),
				),
			}, {
				Config: testAccPolicySettingStringConfig(stringArrayPolicyType, "<<EOF\n- b\n- a\n- d\nEOF", "REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", fmt.Sprintf("%v", []string{"b", "a", "d"})),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value_source", "- b\n- a\n- d\n"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "REQUIRED"),
				),
			},
		},
	})
}

func TestAccPolicySetting_ArrayEncrypted(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicySettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicySettingStringConfigWithPgp(stringArrayPolicyType, "<<EOF\n- a\n- b\n- c\nEOF", "REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "REQUIRED"),
				),
			}, {
				Config: testAccPolicySettingStringConfigWithPgp(stringArrayPolicyType, "<<EOF\n- b\n- a\n- d\nEOF", "REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "REQUIRED"),
				),
			},
		},
	})
}

func TestAccPolicySetting_SecretUnencrypted(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicySettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicySettingStringConfig(secretPolicyType, "test1", "REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "value", "test1"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "REQUIRED"),
				),
			},
		},
	})
}

func TestAccPolicySetting_SecretEncrypted(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicySettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicySettingStringConfigWithPgp(secretPolicyType, "test1", "REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "REQUIRED"),
				),
			},
		},
	})
}

func TestAccPolicySetting_Precedence_Value_Check(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicySettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicySettingPrecedenceAttr("REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "REQUIRED"),
				),
			},
		},
	})
}

func TestAccPolicySetting_Null_Value_Check(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicySettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicySettingStringConfig(stringPolicyType, " ", "REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicySettingExists("turbot_policy_setting.test_policy"),
					resource.TestCheckResourceAttr(
						"turbot_policy_setting.test_policy", "precedence", "REQUIRED"),
				),
			},
		},
	})
}

func testAccPolicySettingPrecedenceAttr(value string) string {
	return fmt.Sprintf(`
resource "turbot_folder" "test" {
  parent = "tmod:@turbot/turbot#/"
  title = "provider_test"
  description = "test folder"
}
resource "turbot_policy_setting" "test_policy" {
  resource    = "${turbot_folder.test.id}"
  type = "tmod:@turbot/aws-s3#/policy/types/bucketTagsTemplate"
   value = <<EOF
  {
    "test2": %s,
    "automation2": "no"
  }
  EOF
}
`, value)
}

// configs
var stringPolicyType = "tmod:@turbot/provider-policy-test#/policy/types/stringPolicy"
var intPolicyType = "tmod:@turbot/provider-policy-test#/policy/types/integerPolicy"
var stringArrayPolicyType = "tmod:@turbot/provider-policy-test#/policy/types/stringArrayPolicy"
var secretPolicyType = "tmod:@turbot/provider-policy-test#/policy/types/secretPolicy"
var stringPolicyTemplate = "{% if $.account.Id == '650022101893' %}Skip{% else %}'Check: Configured'{% endif %}"
var stringPolicyTemplateInput = "{ account{ Id } }"

func testAccPolicySettingStringConfig(policyType, value string, precedence string) string {
	// wrap value in quotes if necessary (this is so we can support heredoc)
	if !strings.HasPrefix(value, `<<EOF`) {
		value = fmt.Sprintf(`"%s"`, value)
	}

	return buildConfig(policyType, value, precedence)
}

func testAccPolicySettingStringConfigWithPgp(policyType, value string, precedence string) string {
	// wrap value in quotes if necessary (this is so we can support heredoc)
	if !strings.HasPrefix(value, `<<EOF`) {
		value = fmt.Sprintf(`"%s"`, value)
	}

	config := fmt.Sprintf(`
resource "turbot_policy_setting" "test_policy" {
	resource = "tmod:@turbot/turbot#/"
	type = "%s"
	value = %s
	precedence = "%s"
	pgp_key = "mQINBF0oedsBEADVfMPaCVRwfaBar8PliWUKU/Q85EiECnfAcfsLyH9TM47o3lhYdH+CkNUvv/1Qqo43ScyGyMRkgw0beQb4jKNdQeSvmsEXl+X9WCHvo2X2fkElaCy74qkilfODmML3Cb7cW6R9j4p0LgF4I42KX5wLqJQy0WV+da4iGuFaJIqDjRjG8a35jxF8cYhLgeh31lSA+ekXTN1e33Ni7ZR1AMBhzJdUjjRZlERPikPmswYO1eiQvRC2viW+Sy22L//ujMJwyAL+tUo5CgJTgf2DXJgf7wOjYJF3fcj9thcxjNRDvs/36xpGESZ3LxPLpP1KClHUfZ0ulrDenH85CiGADezzlF58i6BYWeegEdMcTeVVcbZGhCSgmoX6YjnG31NqdUIayQP3hrk0AGRkySvzcFxP59+PP/jdYnSjNp9hNqB/qpk4isyyvu51B1wnsp6aRAdQA7IrjfF9Q2quVddMO/a2ticAxfUKfWpjYtiucHmFolNmjxQW9S8HEPcf/v/8skdf0w5yUymwkY1AVg4ElfBZ77WKrnhEKduxcHCry4a29XfoXq5RyXaxZ7XM0p+M3aPrPQnf6wCEihrDBwie0ETau32sU/a1PRqKeLBvLRsReioO7Ktb0inxWI92JDqV0XU1ldiRnZDTVtEKvUGHb4SWEkG8mnP6L9N9PL9s4TN2LQARAQABtB1LYWkgRGFndWVycmUgPGthaUB0dXJib3QuY29tPokCTgQTAQgAOBYhBB6Aw6cDXQk0Wipz3PT2jeVANv9mBQJdKHnbAhsDBQsJCAcCBhUKCQgLAgQWAgMBAh4BAheAAAoJEPT2jeVANv9mIEoP/1VLVl/S/s+X8YpF7t0iJ9eK8O13uWzAjZDhshLT3QF6ME5NMGWZb4bIFWvPkmbDpD6WIfe3pHo+v3BDcldK5MVZP3b/Tht+fDkBPpRSZisRDnhaflj954Sq9proP6V9Vnh+bZ5QQnUsyQYHuy5o6cn3q+3lKnBIiTWwNgUlUsMv/WKojzN3gaICptGbBHLdIsU7X8TKmhnhGG5e5tcyOPIh11xoQgt0jkYG5ZtvcEsy8mIxkdFu4hajBUvn1KenZjNZJlsb6B/d+kyeRXGmGcNjpt0+61npiouSJQBigHj6zEF7AuBEFY00mvdnb7fFxVovW44OwvGF3SHXQ3sxPlcxDWQxyqsTnfSeSlX+ZqNtioT5Of4cNYUtOUeXX25KkluqhlCuxG2eXHz2ydj/GW0STkPMeRk6SSlK1a7v2kxNI654cLoKZSTFNocXbzoM2xjdmBZvcfAYudYVuV6dvj5njr0lKjmWeJd4ko6L6O+ypBB9Lp3dHkQ+Q+qcIDMdvnOHjt5MQyOR0XuuZmqUCXBheewqYWXSKNAefM5uVwtwkbdd73cC2rrvb+Pe9Xiz2VwkeqZ9EP0h99Rc2RabASVXpcBs/ISSzPl0viB3Hd0e7F23h/OIshl2qm9m2uCJNK4licm0RYqsO/lkHBLbC9Z2loJMLwHsH4TX5pzoEZOduQINBF0oedsBEAC0kY3NT//sbABiY9RhWjU/HOaBw9dikBP/r4uraqei2dVLfVAnUKSof9F1VpySVZgH4mw+cW5Efev+CNEb5TX/pq4IvdLRogQKiSSlWI0T767wjvnwAtxh5sxjWYQIsRJmfVU9fxlSOPI/0vAuvToGnwGaxbn/to3C14z6G4yo+4tgrRnD6OhM+u4lQVXaKx2gE8B1hze+fgDEVnYEMCCpQoaJbH9cyNmZx6oEfMIKpBhkXJ9Y7OBLpQmckU9PcBwQmlfsIyGxr6eRBu7iPzyPD+QjelM23fyUjXoGrVY1iV32WkxFpFK0S5MDRUXL0fJpErj3Cfqod8DH/3MAt1Rj/IF/RkSqhEcmyUP0j8kG5UVwPN9ZZQBnkgnnI49ggucYwUfyeBXvoC61B3n9BtNwc5Ur63nMKZBgwPfVXtCRstdvoZysxlbI/sumymXpPpM6Sa+nCmqHCbran0sYgHWd26nTubrtTHguq6/Abyd0M4crpbwww3FQTJWnrXCbPsPwvQ35Fk5zDvjeEWK5t3PBDWuUq5AKmdJkpfdgRQQ21Lz7UEvGwTf2I8E3r6YS+Hc+kgCC+qyo6bKQ4q3Fo38OkuK+D0d2e268fymwcACEGqODC7frAm6gmvwe5PkoTnrbzGa2u/sV55JCs/Jqo/bSEzKthtu/A2bHjGWaVF7XgwARAQABiQI2BBgBCAAgFiEEHoDDpwNdCTRaKnPc9PaN5UA2/2YFAl0oedsCGwwACgkQ9PaN5UA2/2Z8qxAAnMEDN72h/qytqxtTCRGrjpydtp/Y4s7yuq+yp7A5Jo0h7h6uW+Opv+tX9Y5CyHjTGbFB/aanGiOJvhXFTEUtc1GGuYZv9mvZrH4DVbJa7yTnV7YjOWqaskRSafC4ftNWXdjr2psuhWCtULgeglR3IUQUzQLHq+GGPINZ92XYPaB2Slgd+/HHbbN/cPObqpb8FQYB2ZuDPif/HLnIAsVsfZhPCC23AySc1kQfXVxdblgEL6L85LTfaF8aKxpdX5YHS+imp8ISj3otzDAQWAL3R5m2/KK4bvWFOTOclbiz73wuJ0l0sM6VK+66R9dCPCl8dcIw33BdIBNPFTtUUJyp33tE8EIbJUTOe2OTLxFEMrxWXKf4iIK+AJenyrbKm0lveAAEh3ynCs3Q8zZpi6L9HmDkhWRlh3toLu9Fz0TmXsF63bSJJtgL7yVwu1KWVuE2ZA2s2nch8AaM8Ozr2pZFLjWtcYyboU2Gp5sCO8iGs3QPxw6W+cxCKLJB9W13sFWZMEfSnu9c7tY/X8LkTqUdk74RXzbJL1jl0GztmUw8n1a5MQBnsenQP/HeyR8qYvQ+Uc8o9blEBLPp/CwtTf3xTqARBB8mj7bZ0YVF2Q6s9TKiYdwW1LgbHlvSdHIqetZHhRE+dgRPwrsTeTjJiaYXIU3UX3oTn4wUiEjvN25dhRQ="
}`,
		policyType, value, precedence)
	return config
}

func testAccPolicySettingIntConfig(policyType string, value int, precedence string) string {
	return buildConfig(policyType, fmt.Sprintf("%d", value), precedence)
}

func testAccPolicySettingTemplateConfig(policyType, template, templateInput, precedence string) string {
	return fmt.Sprintf(`
resource "turbot_policy_setting" "test_policy" {
	resource = "tmod:@turbot/turbot#/"
	type = "%s"
	template = "%s"
	template_input = "%s"
	precedence = "%s"
}`, policyType, template, templateInput, precedence)
}

func buildConfig(policyType, value string, precedence string) string {

	config := fmt.Sprintf(`
resource "turbot_policy_setting" "test_policy" {
	resource = "tmod:@turbot/turbot#/"
	type = "%s"
	value = %s
	precedence = "%s"
}`,
		policyType, value, precedence)
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
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadPolicySetting(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicySettingDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_policy_setting" {
			_, err := client.ReadPolicySetting(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Alert still exists")
			}
			if !apiClient.NotFoundError(err) {
				return fmt.Errorf("expected 'not found' error, got %s", err)
			}
		}
	}

	return nil
}
