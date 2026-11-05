package turbot

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/machinebox/graphql"
	"github.com/turbot/terraform-provider-turbot/apiClient"
)

func TestAccTypedFilter(t *testing.T) {
	accountName := "turbot_typed_filter.account"
	resourceName := "turbot_typed_filter.resource"
	var originalAccountId string
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTypedFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTypedFilterConfig("tf-acc-account-filter", "title:tf-acc-no-such-account", "$.Name:tf-acc-no-such-bucket", "tf-acc-resource-filter"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTypedFilterExists(accountName, "Account"),
					testAccCheckTypedFilterExists(resourceName, "Resource"),
					resource.TestCheckResourceAttr(accountName, "title", "tf-acc-account-filter"),
					resource.TestCheckResourceAttr(accountName, "type", "Account"),
					resource.TestCheckResourceAttr(accountName, "filter", "title:tf-acc-no-such-account"),
					resource.TestCheckResourceAttr(accountName, "description", "acceptance test account filter"),
					resource.TestCheckResourceAttr(accountName, "tags.keep", "a"),
					resource.TestCheckResourceAttr(accountName, "tags.drop", "b"),
					resource.TestCheckResourceAttr(resourceName, "type", "Resource"),
					resource.TestCheckResourceAttr(resourceName, "filter", "$.Name:tf-acc-no-such-bucket"),
					testAccRememberId(accountName, &originalAccountId),
				),
			},
			{
				Config: testAccTypedFilterConfigWithTags("tf-acc-account-filter-renamed", "title:tf-acc-no-such-account-2", "$.Name:tf-acc-no-such-bucket-2", "tf-acc-resource-filter-renamed", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTypedFilterExists(accountName, "Account"),
					testAccCheckTypedFilterExists(resourceName, "Resource"),
					resource.TestCheckNoResourceAttr(accountName, "tags.drop"),
					resource.TestCheckResourceAttr(accountName, "tags.keep", "a"),
					resource.TestCheckResourceAttr(accountName, "title", "tf-acc-account-filter-renamed"),
					resource.TestCheckResourceAttr(accountName, "filter", "title:tf-acc-no-such-account-2"),
					resource.TestCheckResourceAttr(resourceName, "title", "tf-acc-resource-filter-renamed"),
					resource.TestCheckResourceAttr(resourceName, "filter", "$.Name:tf-acc-no-such-bucket-2"),
					testAccCheckSameId(accountName, &originalAccountId),
				),
			},
			{
				ResourceName:      accountName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() { testAccDeleteTypedFilter(t, originalAccountId) },
				Config:    testAccTypedFilterConfigWithTags("tf-acc-account-filter-renamed", "title:tf-acc-no-such-account-2", "$.Name:tf-acc-no-such-bucket-2", "tf-acc-resource-filter-renamed", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTypedFilterExists(accountName, "Account"),
					testAccCheckDifferentId(accountName, &originalAccountId),
				),
			},
		},
	})
}

func testAccTypedFilterConfig(accountTitle, accountFilter, resourceFilter, resourceTitle string) string {
	return testAccTypedFilterConfigWithTags(accountTitle, accountFilter, resourceFilter, resourceTitle, "\t\tdrop = \"b\"\n")
}

func testAccTypedFilterConfigWithTags(accountTitle, accountFilter, resourceFilter, resourceTitle, extraTags string) string {
	return fmt.Sprintf(`
resource "turbot_typed_filter" "account" {
	title       = "%s"
	type        = "Account"
	filter      = "%s"
	description = "acceptance test account filter"
	tags = {
		keep = "a"
%s	}
}

resource "turbot_typed_filter" "resource" {
	title  = "%s"
	type   = "Resource"
	filter = "%s"
}
`, accountTitle, accountFilter, extraTags, resourceTitle, resourceFilter)
}

func testAccCheckTypedFilterExists(name, wantType string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		typedFilter, err := client.ReadTypedFilter(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching %s: %s", name, err)
		}
		if typedFilter.Type != wantType {
			return fmt.Errorf("%s has type %s, want %s", name, typedFilter.Type, wantType)
		}
		return nil
	}
}

func testAccRememberId(name string, id *string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		*id = state.RootModule().Resources[name].Primary.ID
		return nil
	}
}

func testAccCheckSameId(name string, id *string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		if got := state.RootModule().Resources[name].Primary.ID; got != *id {
			return fmt.Errorf("%s was replaced: id %s, want %s", name, got, *id)
		}
		return nil
	}
}

func testAccCheckDifferentId(name string, id *string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		if got := state.RootModule().Resources[name].Primary.ID; got == *id {
			return fmt.Errorf("%s still has the deleted id %s: refresh did not drop it from state", name, got)
		}
		return nil
	}
}

func testAccDeleteTypedFilter(t *testing.T, id string) {
	client := testAccProvider.Meta().(*apiClient.Client)
	if err := client.DeleteTypedFilter(id); err != nil {
		t.Fatalf("deleting typed filter %s outside Terraform: %s", id, err)
	}
}

func testAccCheckTypedFilterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "turbot_typed_filter" {
			continue
		}
		_, err := client.ReadTypedFilter(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("typed filter %s still exists", rs.Primary.ID)
		}
		if !apiClient.IsTypedFilterNotFound(err) {
			return fmt.Errorf("expected typed filter %s to be not found, got %s", rs.Primary.ID, err)
		}
	}
	return nil
}

func typedFilterTestClient(t *testing.T, body string) *apiClient.Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, body)
	}))
	t.Cleanup(server.Close)
	return &apiClient.Client{Graphql: graphql.NewClient(server.URL + "/graphql")}
}

func typedFilterTestData(t *testing.T) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, resourceTurbotTypedFilter().Schema, map[string]interface{}{
		"title":  "tf-acc-x",
		"type":   "Account",
		"filter": "title:x",
	})
	d.SetId("401101111592184")
	return d
}

func TestTypedFilterReadKeepsStateOnOtherError(t *testing.T) {
	client := typedFilterTestClient(t, `{"errors":[{"message":"Forbidden: Insufficient permissions"}]}`)
	d := typedFilterTestData(t)

	err := resourceTurbotTypedFilterRead(d, client)

	if err == nil {
		t.Fatal("expected the read error to fail the refresh")
	}
	if d.Id() != "401101111592184" {
		t.Fatalf("a non-not-found error must keep the filter in state, id is %q", d.Id())
	}
}

func TestTypedFilterReadDropsStateWhenGone(t *testing.T) {
	client := typedFilterTestClient(t, `{"data":{"typedFilter":null}}`)
	d := typedFilterTestData(t)

	if err := resourceTurbotTypedFilterRead(d, client); err != nil {
		t.Fatalf("a missing filter is not a refresh error: %s", err)
	}
	if d.Id() != "" {
		t.Fatalf("a deleted filter must drop from state, id is %q", d.Id())
	}
}

func TestTypedFilterTypeValidation(t *testing.T) {
	for _, ok := range []string{"Account", "Resource"} {
		if _, errs := validateTypedFilterType(ok, "type"); len(errs) != 0 {
			t.Errorf("%s must be accepted: %v", ok, errs)
		}
	}
	if _, errs := validateTypedFilterType("account", "type"); len(errs) == 0 {
		t.Error("lower-case type must be rejected")
	}
}
