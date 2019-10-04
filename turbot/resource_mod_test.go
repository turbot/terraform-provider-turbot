package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"testing"
)

// test suites

func TestAccMod(t *testing.T) {
	latestProviderTestVersion := "5.0.15"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccModDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMod_v5_0_0_Config(),
				Check: resource.ComposeTestCheckFunc(
					testAccModExists("turbot_mod.test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "mod", "provider-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version_current", "5.0.0"),
				),
			},
			{
				Config: testAccMod_ge_v5_0_0_Config(),
				Check: resource.ComposeTestCheckFunc(
					testAccModExists("turbot_mod.test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "mod", "provider-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version_current", latestProviderTestVersion),
				),
			},
			{
				Config: testAccMod_v5_0_1_Config(),
				Check: resource.ComposeTestCheckFunc(
					testAccModExists("turbot_mod.test"),
					testAccModExists("turbot_mod.test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "mod", "provider-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version_current", "5.0.1"),
				),
			},
			{
				Config: testAccModWildCardConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccModExists("turbot_mod.test"),
					testAccModExists("turbot_mod.test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "mod", "provider-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version_current", latestProviderTestVersion),
				),
			},
			{
				Config: testAccModWildCardConfig2(),
				Check: resource.ComposeTestCheckFunc(
					testAccModExists("turbot_mod.test"),
					testAccModExists("turbot_mod.test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "mod", "provider-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version_current", latestProviderTestVersion),
				),
			},
			{
				Config: testAccMod_lt_v5_0_3_Config(),
				Check: resource.ComposeTestCheckFunc(
					testAccModExists("turbot_mod.test"),
					testAccModExists("turbot_mod.test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "mod", "provider-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version_current", "5.0.2"),
				),
			},
			{
				Config: testAccModNoVersionConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccModExists("turbot_mod.test"),
					testAccModExists("turbot_mod.test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "mod", "provider-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version_current", latestProviderTestVersion),
				),
			},
		},
	})
}

// configs
func testAccMod_v5_0_0_Config() string {
	return `
resource "turbot_mod" "test" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "provider-test"
	version = "5.0.0"
}
`
}

func testAccMod_v5_0_1_Config() string {
	return `
resource "turbot_mod" "test" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "provider-test"
	version = "5.0.1"
}
`
}

func testAccMod_ge_v5_0_0_Config() string {
	return `
resource "turbot_mod" "test" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "provider-test"
	version = ">=5.0.0"
}
`
}

func testAccMod_lt_v5_0_3_Config() string {
	return `
resource "turbot_mod" "test" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "provider-test"
	version = "<5.0.3"
}
`
}

func testAccModWildCardConfig() string {
	return `
resource "turbot_mod" "test" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "provider-test"
	version = "*"
}
`
}

func testAccModWildCardConfig2() string {
	return `
resource "turbot_mod" "test" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "provider-test"
	version = "5.0.*"
}
`
}

func testAccModNoVersionConfig() string {
	return `
resource "turbot_mod" "test" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "provider-test"
}
`
}

// helper functions
func testAccModExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiclient.Client)
		_, err := client.ReadMod(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccModDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mod" {
			continue
		}
		_, err := client.ReadMod(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		if !apiclient.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}

	return nil
}
