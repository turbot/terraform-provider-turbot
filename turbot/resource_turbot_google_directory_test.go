package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
	"testing"
)

func TestAccGoogleDirectory_Pgp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleDirectoryConfigPgp(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleDirectoryExists("turbot_google_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "title", "google_directory_test_provider"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "description", "test directory"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "profile_id_template", "profileemail"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "client_id", "provider-test.apps.google.com"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "parent", "tmod:@turbot/turbot#/"),
				),
			},
		},
	})
}

// test suites
func TestAccGoogleDirectory_Basic(t *testing.T) {
	resourceName := "turbot_google_directory.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleDirectoryConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleDirectoryExists("turbot_google_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "title", "google_directory_test_provider"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "description", "test directory"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_secret"},
			},
			{
				Config: testAccGoogleDirectoryUpdateTitleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleDirectoryExists("turbot_google_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "title", "google_directory_test_provider2"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "description", "test directory"),
				),
			},
			{
				Config: testAccGoogleDirectoryUpdateDescConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleDirectoryExists("turbot_google_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "title", "google_directory_test_provider"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "description", "test directory for turbot terraform provider"),
				),
			},
		},
	})
}

func TestAccGoogleDirectory_tags(t *testing.T) {
	resourceName := "turbot_google_directory.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleDirectoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleDirectoryTagsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleDirectoryExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "hosted_name", "turbot.com"),
					resource.TestCheckResourceAttr(resourceName, "tags.tag1", "tag1value"),
					resource.TestCheckResourceAttr(resourceName, "tags.tag2", "tag2value"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_secret"},
			},
		},
	})
}

// configs
func testAccGoogleDirectoryConfigPgp() string {
	return `
resource "turbot_google_directory" "test" {
	title = "google_directory_test_provider"
	profile_id_template = "profileemail"
	client_id = "provider-test.apps.google.com"
	client_secret = "scqXnRczuyve329"
	parent = "tmod:@turbot/turbot#/"
	description = "test directory"
	pgp_key = "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"
}
`
}
func testAccGoogleDirectoryConfig() string {
	return `
resource "turbot_google_directory" "test" {
	title = "google_directory_test_provider"
	profile_id_template = "profileemail"
	client_id = "provider-test.apps.google.com"
	client_secret = "scqXnRczuyve329"
	parent = "tmod:@turbot/turbot#/"
	description = "test directory"
}
`
}

func testAccGoogleDirectoryUpdateTitleConfig() string {
	return `
resource "turbot_google_directory" "test" {
	title = "google_directory_test_provider2"
	profile_id_template = "profileemail"
	client_id = "provider-test.apps.google.com"
	client_secret = "scqXnRczuyve329"
	parent = "tmod:@turbot/turbot#/"
	description = "test directory"
}
`
}

func testAccGoogleDirectoryUpdateDescConfig() string {
	return `
resource "turbot_google_directory" "test" {
	title = "google_directory_test_provider"
	profile_id_template = "profileemail"
	client_id = "provider-test.apps.google.com"
	client_secret = "scqXnRczuyve329"
	parent = "tmod:@turbot/turbot#/"
	description = "test directory for turbot terraform provider"
}
`
}

func testAccGoogleDirectoryTagsConfig() string {
	return `
resource "turbot_google_directory" "test" {
	title = "google_directory_test_provider"
	profile_id_template = "profileemail"
	client_id = "provider-test.apps.google.com"
	hosted_name = "turbot.com"
	client_secret = "scqXnRczuyve329"
	parent = "tmod:@turbot/turbot#/"
	description = "test directory for turbot terraform provider"
	tags = {
		tag1 = "tag1value"
		tag2 = "tag2value"
	}
}`
}

// helper functions
func testAccCheckGoogleDirectoryExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadGoogleDirectory(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckGoogleDirectoryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_google_directory" {
			_, err := client.ReadGoogleDirectory(rs.Primary.ID)
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
