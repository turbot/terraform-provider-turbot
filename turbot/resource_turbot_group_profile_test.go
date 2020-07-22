package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"testing"
)

// test suites
func TestAccGroupProfile_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGroupProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAcGroupProfileConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupProfileExists("turbot_group_profile.test"),
					resource.TestCheckResourceAttr("turbot_group_profile.test", "title", "Snape"),
					resource.TestCheckResourceAttr("turbot_group_profile.test", "status", "Active"),
				),
			},
		},
	})
}

// configs
func testAcGroupProfileConfig() string {
	return `
resource "turbot_group_profile" "test" {
  	directory = "7878778"
  	title  = "test"
  	status = "ACTIVE"
  	group_profile_id = "7879877"
}
`
}

// helper functions
func testAccCheckGroupProfileExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no reecord id is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadGroupProfile(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckGroupProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_group_profile" {
			_, err := client.ReadGroupProfile(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("alert still exists")
			}
			if !apiClient.NotFoundError(err) {
				return fmt.Errorf("expected 'not found' error, got %s", err)
			}
		}
	}

	return nil
}
