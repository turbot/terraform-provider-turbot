package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"testing"
)

// test suites
func TestAccGoogleDirectory(t *testing.T) {
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
						"turbot_google_directory.test", "description", "test Directory"),
				),
			},
			{
				Config: testAccGoogleDirectoryUpdateTitleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleDirectoryExists("turbot_google_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "title", "google_directory_test_provider2"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "description", "test Directory"),
				),
			},
			{
				Config: testAccGoogleDirectoryUpdateDescConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleDirectoryExists("turbot_google_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "title", "google_directory_test_provider"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "description", "test Directory for turbot terraform provider"),
				),
			}, {
				Config: testAccGoogleDirectoryTagsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleDirectoryExists("turbot_google_directory.test"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "title", "google_directory_test_provider"),
					resource.TestCheckResourceAttr(
						"turbot_google_directory.test", "description", "test Directory for turbot terraform provider"),
					resource.TestCheckResourceAttr(
						"turbot_folder.test", "tags.Name", "tags test"),
					resource.TestCheckResourceAttr(
						"turbot_folder.test", "tags.Environment", "foo"),
				),
			},
		},
	})
}

func testAccGoogleDirectoryConfig() string {
	return `
resource "turbot_google_directory" "test" {
	title = "google_directory_test_provider"
	profile_id_template = "profileemail"
	status = "New"
	directory_type = "google"
	client_id = "GoogleDirTest4"
	client_secret = "fb-tbevaACsBKQHthzba-PH9"
	parent = "tmod:@turbot/turbot#/"
	description = "test Directory"
}
`
}

func testAccGoogleDirectoryUpdateTitleConfig() string {
	return `
resource "turbot_google_directory" "test" {
	title = "google_directory_test_provider2"
	profile_id_template = "profileemail"
	status = "New"
	directory_type = "google"
	client_id = "GoogleDirTest4"
	client_secret = "fb-tbevaACsBKQHthzba-PH9"
	parent = "tmod:@turbot/turbot#/"
	description = "test Directory"
}
`
}

func testAccGoogleDirectoryUpdateDescConfig() string {
	return `
resource "turbot_google_directory" "test" {
	title = "google_directory_test_provider"
	profile_id_template = "profileemail"
	status = "New"
	directory_type = "google"
	client_id = "GoogleDirTest4"
	client_secret = "fb-tbevaACsBKQHthzba-PH9"
	parent = "tmod:@turbot/turbot#/"
	description = "test Directory for turbot terraform provider"
}
`
}

func testAccGoogleDirectoryTagsConfig() string {
	return `
resource "turbot_google_directory" "test" {
	title = "google_directory_test_provider"
	profile_id_template = "profileemail"
	status = "New"
	directory_type = "google"
	client_id = "GoogleDirTest4"
	client_secret = "fb-tbevaACsBKQHthzba-PH9"
	parent = "tmod:@turbot/turbot#/"
	description = "test Directory for turbot terraform provider"
	tags = {
		  "Name" = "tags test"
		  "Environment" = "foo"
	}
}
`
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
		client := testAccProvider.Meta().(*apiclient.Client)
		_, err := client.ReadGoogleDirectory(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckGoogleDirectoryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "Directory" {
			continue
		}
		_, err := client.ReadGoogleDirectory(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		if !apiclient.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}

	return nil
}
