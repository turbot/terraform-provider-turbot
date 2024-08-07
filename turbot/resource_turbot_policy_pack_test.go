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
func TestAccPolicyPack_Basic(t *testing.T) {
	resourceName := "turbot_policy_pack.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyPackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyPackConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyPackExists("turbot_policy_pack.test"),
					resource.TestCheckResourceAttr("turbot_policy_pack.test", "title", "policy_pack"),
					resource.TestCheckResourceAttr("turbot_policy_pack.test", "description", "Policy Pack Testing"),
					resource.TestCheckResourceAttr("turbot_policy_pack.test", "parent", "178806508050433"),
					resource.TestCheckResourceAttr("turbot_policy_pack.test", "filter", "resourceType:181381985925765 $.turbot.tags.a:b"),
				),
			},
			{
				Config: testAccPolicyPackUpdateDescConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyPackExists("turbot_policy_pack.test"),
					resource.TestCheckResourceAttr("turbot_policy_pack.test", "title", "policy_pack"),
					resource.TestCheckResourceAttr("turbot_policy_pack.test", "description", "Policy Pack updated"),
					resource.TestCheckResourceAttr("turbot_policy_pack.test", "parent", "178806508050433"),
					resource.TestCheckResourceAttr("turbot_policy_pack.test", "filter", "resourceType:181381985925765 $.turbot.tags.a:b"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"filter"},
			},
		},
	})
}

// configs
func testAccPolicyPackConfig() string {
	return `
resource "turbot_policy_pack" "test" {
	filter = "resourceType:181381985925765 $.turbot.tags.a:b"
	description = "Policy Pack Testing"
	title = "policy_pack"
}
`
}

func testAccPolicyPackUpdateDescConfig() string {
	return `
resource "turbot_policy_pack" "test" {
	filter = "resourceType:181381985925765 $.turbot.tags.a:b"
	description = "Policy Pack updated"
	title ="policy_pack"
}
`
}

// helper functions
func testAccCheckPolicyPackExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadSmartFolder(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckPolicyPackDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_policy_pack" {
			_, err := client.ReadSmartFolder(rs.Primary.ID)
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
