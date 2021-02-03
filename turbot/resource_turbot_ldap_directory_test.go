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
func TestAccLdapDirectory_Basic(t *testing.T) {
	resourceName := "turbot_ldap_directory.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLdapDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLdapDirectoryConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLdapDirectoryExists("turbot_ldap_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_ldap_directory.test", "title", "Microsoft LDAP dir"),
					resource.TestCheckResourceAttr(
						"turbot_ldap_directory.test", "distinguished_name", "CN=Turbot"),
				),
			},
			{
				Config: testAccLdapDirectoryUpdateUrlConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLdapDirectoryExists("turbot_ldap_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_ldap_directory.test", "title", "Microsoft LDAP dir"),
					resource.TestCheckResourceAttr(
						"turbot_ldap_directory.test", "url", "xws"),
				),
			},
			{
				Config: testAccLdapDirectoryUpdateTitleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLdapDirectoryExists("turbot_ldap_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_ldap_directory.test", "title", "Azure LDAP dir"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccLdapDirectory_tags(t *testing.T) {
	resourceName := "turbot_ldap_directory.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLdapDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLdapDirectoryTagsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLdapDirectoryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.tag1", "tag1value"),
					resource.TestCheckResourceAttr(resourceName, "tags.tag2", "tag2value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

// configs
func testAccLdapDirectoryConfig() string {
	return `
resource "turbot_ldap_directory" "test" {
	parent = "tmod:@turbot/turbot#/"
  	title = "Microsoft LDAP dir"
  	profile_id_template =  "{{profile.email}}"
  	distinguished_name = "CN=Turbot"
  	password = "x7hjFeErf0_+"
  	url = "xw"
  	base = "xw"
  	tls_enabled = false
  	reject_unauthorized = false
    tags = {
		  tag1 = "tag1value"
		  tag2 = "tag2value"
	}
}
`
}

func testAccLdapDirectoryUpdateUrlConfig() string {
	return `
resource "turbot_ldap_directory" "test" {
	parent = "tmod:@turbot/turbot#/"
  	title = "Microsoft LDAP dir"
  	profile_id_template =  "{{profile.email}}"
  	distinguished_name = "CN=Turbot"
  	password = "x7hjFeErf0_+"
  	url = "xws"
  	base = "xw"
  	tls_enabled = false
  	reject_unauthorized = false
    tags = {
		  tag1 = "tag1value"
		  tag2 = "tag2value"
	}
}
`
}

func testAccLdapDirectoryUpdateTitleConfig() string {
	return `
resource "turbot_ldap_directory" "test" {
	parent = "tmod:@turbot/turbot#/"
  	title = "Azure LDAP dir"
  	profile_id_template =  "{{profile.email}}"
  	distinguished_name = "CN=Turbot"
  	password = "x7hjFeErf0_+"
  	url = "xw"
  	base = "xw"
  	tls_enabled = false
  	reject_unauthorized = false
    tags = {
		  tag1 = "tag1value"
		  tag2 = "tag2value"
	}
}
`
}

func testAccLdapDirectoryTagsConfig() string {
	return `
resource "turbot_ldap_directory" "test" {
	parent = "tmod:@turbot/turbot#/"
  	title = "Microsoft LDAP dir"
  	profile_id_template =  "{{profile.email}}"
  	distinguished_name = "CN=Turbot"
  	password = "x7hjFeErf0_+"
  	url = "xw"
  	base = "xw"
  	tls_enabled = false
  	reject_unauthorized = false
    tags = {
		  tag1 = "tag1value"
		  tag2 = "tag2value"
	}
}
`
}

// helper functions
func testAccCheckLdapDirectoryExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadLdapDirectory(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckLdapDirectoryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_ldap_directory" {
			_, err := client.ReadLdapDirectory(rs.Primary.ID)
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
