package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
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
						"turbot_ldap_directory.test", "title", "provider_test"),
					resource.TestCheckResourceAttr(
						"turbot_ldap_directory.test", "description", "test directory"),
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

//func TestAccLdapDirectory_tags(t *testing.T) {
//	resourceName := "turbot_ldap_directory.test"
//	resource.Test(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers:    testAccProviders,
//		CheckDestroy: testAccCheckLdapDirectoryDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccDirectoryTagsConfig(),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckLdapDirectoryExists(resourceName),
//					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
//					resource.TestCheckResourceAttr(resourceName, "tags.tag1", "tag1value"),
//					resource.TestCheckResourceAttr(resourceName, "tags.tag2", "tag2value"),
//				),
//			},
//			{
//				ResourceName:      resourceName,
//				ImportState:       true,
//				ImportStateVerify: true,
//			},
//		},
//	})
//}

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
}
`
}

/*func testAccLdapDirectoryUpdateDescConfig() string {
	return `
resource "turbot_ldap_directory" "test" {
	parent = "tmod:@turbot/turbot#/"
	title = "provider_test"
	description = "test directory for turbot terraform provider"
	profile_id_template = "{{profile.email}}"
}
`
}*/

//func testAccDirectoryUpdateTitleConfig() string {
//	return `
//resource "turbot_ldap_directory" "test" {
//	parent = "tmod:@turbot/turbot#/"
//	title = "provider_test_refactor"
//	description = "test directory"
//	profile_id_template = "{{profile.email}}"
//}
//`
//}

//func testAccDirectoryTagsConfig() string {
//	return `
//resource "turbot_ldap_directory" "test" {
//	parent = "tmod:@turbot/turbot#/"
//	title = "provider_test_refactor"
//	description = "test directory"
//	profile_id_template = "{{profile.email}}"
//	tags = {
//		  tag1 = "tag1value"
//		  tag2 = "tag2value"
//	}
//}
//`
//}

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
			if !apiClient.NotFoundError(err) {
				return fmt.Errorf("expected 'not found' error, got %s", err)
			}
		}
	}

	return nil
}
