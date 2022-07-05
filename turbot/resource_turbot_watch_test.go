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
func TestAccWatch_Basic(t *testing.T) {
	resourceName := "turbot_watch.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWatchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWatchConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFolderExists("turbot_folder.test"),
					testAccCheckWatchExists("turbot_watch.test"),
					resource.TestCheckResourceAttr("turbot_watch.test", "action", "tmod:@turbot/firehose-aws-sns#/action/types/router"),
				),
			},
			{
				Config: testAccWatchUpdateFilterConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWatchExists("turbot_watch.test"),
					resource.TestCheckResourceAttr("turbot_watch.test", "action", "tmod:@turbot/firehose-aws-sns#/action/types/router"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"filters"},
			},
		},
	})
}

// configs
func testAccWatchConfig() string {
	return `
	resource "turbot_folder" "test" {
		parent = "tmod:@turbot/turbot#/"
		title = "provider_test"
		description = "test folder"
	}

	resource "turbot_watch" "test" {
		resource = turbot_folder.test.id
		action   = "tmod:@turbot/firehose-aws-sns#/action/types/router"
		filters  = ["resourceId:${turbot_folder.test.id} level:self,descendant"]
	}`
}

func testAccWatchUpdateFilterConfig() string {
	return `
	resource "turbot_folder" "test" {
		parent = "tmod:@turbot/turbot#/"
		title = "provider_test"
		description = "test folder"
	}

	resource "turbot_watch" "test" {
		resource = turbot_folder.test.id
		action   = "tmod:@turbot/firehose-aws-sns#/action/types/router"
		filters  = ["resourceId:${turbot_folder.test.id} level:self"]
	}`
}

// helper functions
func testAccCheckWatchExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadWatch(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckWatchDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_watch" {
			_, err := client.ReadWatch(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("alert still exists")
			}
			if !errors.NotFoundError(err) {
				return fmt.Errorf("expected 'not found' error, got %s", err)
			}
		}
	}

	return nil
}
