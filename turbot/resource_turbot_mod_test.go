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

func TestAccMod_Basic(t *testing.T) {
	resourceName := "turbot_mod.test"
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
						"turbot_mod.test", "mod", "turbot-terraform-provider-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version_current", "5.0.0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"version"},
			},
		},
	})
}

func TestAccMod_AutoPatchVersionUpgrade(t *testing.T) {
	latestProviderTestVersion := "5.1.0"
	resourceName := "turbot_mod.test"
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
						"turbot_mod.test", "mod", "turbot-terraform-provider-test"),
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
						"turbot_mod.test", "mod", "turbot-terraform-provider-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version_current", latestProviderTestVersion),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"version"},
			},
		},
	})
}

func TestAccMod_AutoPatchVersionDowngrade(t *testing.T) {
	latestProviderTestVersion := "5.1.0"
	resourceName := "turbot_mod.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccModDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMod_ge_v5_0_0_Config(),
				Check: resource.ComposeTestCheckFunc(
					testAccModExists("turbot_mod.test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "mod", "turbot-terraform-provider-test"),
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
						"turbot_mod.test", "mod", "turbot-terraform-provider-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version_current", "5.0.1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"version"},
			},
		},
	})
}

func TestAccMod_WildCardVersion(t *testing.T) {
	latestProviderTestVersion := "5.1.0"
	resourceName := "turbot_mod.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccModDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModWildCardConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccModExists("turbot_mod.test"),
					testAccModExists("turbot_mod.test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "mod", "turbot-terraform-provider-test"),
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
						"turbot_mod.test", "mod", "turbot-terraform-provider-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version_current", "5.0.2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"version"},
			},
		},
	})
}

func TestAccMod_LowerThanVersion(t *testing.T) {
	resourceName := "turbot_mod.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccModDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMod_lt_v5_0_3_Config(),
				Check: resource.ComposeTestCheckFunc(
					testAccModExists("turbot_mod.test"),
					testAccModExists("turbot_mod.test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "mod", "turbot-terraform-provider-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version_current", "5.0.2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"version"},
			},
		},
	})
}

func TestAccMod_NoVersion(t *testing.T) {
	latestProviderTestVersion := "5.1.0"
	resourceName := "turbot_mod.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccModDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModNoVersionConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccModExists("turbot_mod.test"),
					testAccModExists("turbot_mod.test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "mod", "turbot-terraform-provider-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version_current", latestProviderTestVersion),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"version"},
			},
		},
	})
}

func TestAccMod_NoParentResource(t *testing.T) {
	resourceName := "turbot_mod.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccModDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModNoParentConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccModExists("turbot_mod.test"),
					testAccModExists("turbot_mod.test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "org", "turbot"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "mod", "turbot-terraform-provider-test"),
					resource.TestCheckResourceAttr(
						"turbot_mod.test", "version_current", "5.0.0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"version"},
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
	mod = "turbot-terraform-provider-test"
	version = "5.0.0"
}
`
}

func testAccMod_v5_0_1_Config() string {
	return `
resource "turbot_mod" "test" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "turbot-terraform-provider-test"
	version = "5.0.1"
}
`
}

func testAccMod_ge_v5_0_0_Config() string {
	return `
resource "turbot_mod" "test" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "turbot-terraform-provider-test"
	version = ">=5.0.0"
}
`
}

func testAccMod_lt_v5_0_3_Config() string {
	return `
resource "turbot_mod" "test" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "turbot-terraform-provider-test"
	version = "<5.0.3"
}
`
}

func testAccModWildCardConfig() string {
	return `
resource "turbot_mod" "test" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "turbot-terraform-provider-test"
	version = "*"
}
`
}

func testAccModWildCardConfig2() string {
	return `
resource "turbot_mod" "test" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "turbot-terraform-provider-test"
	version = "5.0.*"
}
`
}

func testAccModNoVersionConfig() string {
	return `
resource "turbot_mod" "test" {
	parent = "tmod:@turbot/turbot#/"
	org = "turbot"
	mod = "turbot-terraform-provider-test"
}
`
}

func testAccModNoParentConfig() string {
	return `
resource "turbot_mod" "test"{
	org = "turbot"
	mod = "turbot-terraform-provider-test"
	version = "5.0.0"
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
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadMod(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccModDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_mod" {
			_, err := client.ReadMod(rs.Primary.ID)
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
