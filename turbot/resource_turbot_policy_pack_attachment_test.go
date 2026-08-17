package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
	"regexp"
	"strconv"
	"testing"
)

// test suites
func TestAccPolicyPackAttachment_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyPackAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyPackAttachmentConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyPackAttachmentExists("turbot_policy_pack_attachment.test"),
				),
			},
		},
	})
}

// configs
func testAccPolicyPackAttachmentConfig() string {
	return `
resource "turbot_folder" "test" {
  parent = "tmod:@turbot/turbot#/"
  title = "provider_test"
  description = "test folder"
}

resource "turbot_policy_pack" "test" {
  filter      = "resourceType:181381985925765 $.turbot.tags.a:b"
  description = "Policy Pack Testing"
  title       = "policy_pack"
}

resource "turbot_policy_pack_attachment" "test" {
  resource    = turbot_folder.test.id
  policy_pack = turbot_policy_pack.test.id
}
`
}

func TestAccPolicyPackAttachment_Aka(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyPackAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyPackAttachmentAkaConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyPackAttachmentExists("turbot_policy_pack_attachment.test_aka"),
				),
			},
		},
	})
}

// configs
func testAccPolicyPackAttachmentAkaConfig() string {
	return `
resource "turbot_folder" "test_aka" {
  parent      = "tmod:@turbot/turbot#/"
  title       = "provider_test_aka"
  description = "test folder for aka attachment"
}

resource "turbot_policy_pack" "test_aka" {
  filter      = "resourceType:181381985925765 $.turbot.tags.a:b"
  description = "Policy Pack AKA Testing"
  title       = "policy_pack_aka"
  akas        = ["test_policy_pack_aka_acceptance"]
}

resource "turbot_policy_pack_attachment" "test_aka" {
  resource    = turbot_folder.test_aka.id
  policy_pack = "test_policy_pack_aka_acceptance"
  depends_on  = [turbot_policy_pack.test_aka]
}
`
}

// helper functions
func testAccCheckPolicyPackAttachmentExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		policyPackId, resource := parsePolicyPackId(rs.Primary.ID)
		// Verify the ATTACHMENT from the resource side, matching what Exists() does. Reading the
		// pack (the old check) only proves the pack exists, and needs a grant on the pack.
		attached, err := client.PolicyPackAttached(resource, policyPackId)
		if err != nil {
			return fmt.Errorf("error fetching attachment for resource %s. %s", resource, err)
		}
		if !attached {
			return fmt.Errorf("policy pack %s is not attached to resource %s", policyPackId, resource)
		}
		return nil
	}
}

// testAccCheckPolicyPackAttachmentDetached asserts the pack is NOT attached to the resource.
// Pairs with the Exists() path: a detached attachment must drop out of state rather than linger.
func testAccCheckPolicyPackAttachmentDetached(policyPack, resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := testAccProvider.Meta().(*apiClient.Client)
		attached, err := client.PolicyPackAttached(resource, policyPack)
		if err != nil {
			return fmt.Errorf("error reading attachments for resource %s. %s", resource, err)
		}
		if attached {
			return fmt.Errorf("policy pack %s is still attached to resource %s", policyPack, resource)
		}
		return nil
	}
}
func testAccCheckPolicyPackAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "policyPack" {
			continue
		}
		_, err := client.ReadSmartFolder(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("alert still exists")
		}
		if !errors.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}
	return nil
}

// TestAccPolicyPackAttachment_IdAkaMatrix covers every combination of identifying the pack and
// the target by numeric id or by aka. All four must attach, and `policy_pack` must always
// settle to the resolved numeric id in state (the attachSmartFolders mutation rejects akas),
// while `policy_pack_akas` carries the akas so suppressIfAkaMatches prevents a perpetual diff.
func TestAccPolicyPackAttachment_IdAkaMatrix(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyPackAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyPackAttachmentMatrixConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyPackAttachmentExists("turbot_policy_pack_attachment.id_id"),
					testAccCheckPolicyPackAttachmentExists("turbot_policy_pack_attachment.aka_id"),
					testAccCheckPolicyPackAttachmentExists("turbot_policy_pack_attachment.id_aka"),
					testAccCheckPolicyPackAttachmentExists("turbot_policy_pack_attachment.aka_aka"),
					// an aka in config must still resolve to the numeric id in state
					resource.TestMatchResourceAttr("turbot_policy_pack_attachment.aka_aka", "policy_pack",
						regexp.MustCompile(`^[0-9]+$`)),
					resource.TestMatchResourceAttr("turbot_policy_pack_attachment.aka_id", "policy_pack",
						regexp.MustCompile(`^[0-9]+$`)),
					// the akas list must be non-empty so the diff suppressor has something to match.
					// Asserting a specific index would pin an ordering the API does not guarantee,
					// and suppressIfAkaMatches scans the whole list rather than one element.
					testAccCheckAkasContain("turbot_policy_pack_attachment.aka_aka", "policy_pack_akas",
						"test_pack_matrix_aka"),
				),
			},
			{
				// re-planning the same config must produce no diff, proving the aka/id
				// suppression survives a round-trip through state
				Config:   testAccPolicyPackAttachmentMatrixConfig(),
				PlanOnly: true,
			},
		},
	})
}

func testAccPolicyPackAttachmentMatrixConfig() string {
	// Each combination needs its OWN target folder: attaching the same pack to the same
	// resource twice is a duplicate attachment, not a second test case.
	return `
resource "turbot_policy_pack" "matrix" {
  filter      = "resourceType:181381985925765 $.turbot.tags.a:b"
  description = "Policy Pack id/aka matrix testing"
  title       = "policy_pack_matrix"
  akas        = ["test_pack_matrix_aka"]
}

resource "turbot_folder" "m1" {
  parent      = "tmod:@turbot/turbot#/"
  title       = "provider_test_matrix_1"
  description = "matrix target: pack id, resource id"
  akas        = ["test_folder_matrix_1"]
}

resource "turbot_folder" "m2" {
  parent      = "tmod:@turbot/turbot#/"
  title       = "provider_test_matrix_2"
  description = "matrix target: pack aka, resource id"
  akas        = ["test_folder_matrix_2"]
}

resource "turbot_folder" "m3" {
  parent      = "tmod:@turbot/turbot#/"
  title       = "provider_test_matrix_3"
  description = "matrix target: pack id, resource aka"
  akas        = ["test_folder_matrix_3"]
}

resource "turbot_folder" "m4" {
  parent      = "tmod:@turbot/turbot#/"
  title       = "provider_test_matrix_4"
  description = "matrix target: pack aka, resource aka"
  akas        = ["test_folder_matrix_4"]
}

resource "turbot_policy_pack_attachment" "id_id" {
  policy_pack = turbot_policy_pack.matrix.id
  resource    = turbot_folder.m1.id
}

resource "turbot_policy_pack_attachment" "aka_id" {
  policy_pack = "test_pack_matrix_aka"
  resource    = turbot_folder.m2.id
  depends_on  = [turbot_policy_pack.matrix]
}

resource "turbot_policy_pack_attachment" "id_aka" {
  policy_pack = turbot_policy_pack.matrix.id
  resource    = "test_folder_matrix_3"
  depends_on  = [turbot_folder.m3]
}

resource "turbot_policy_pack_attachment" "aka_aka" {
  policy_pack = "test_pack_matrix_aka"
  resource    = "test_folder_matrix_4"
  depends_on  = [turbot_policy_pack.matrix, turbot_folder.m4]
}
`
}

// TestAccPolicyPackAttachment_Import verifies the import path, which routes through Read() and
// therefore through the policy pack aka lookup that this fix changed.
func TestAccPolicyPackAttachment_Import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyPackAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyPackAttachmentConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyPackAttachmentExists("turbot_policy_pack_attachment.test"),
				),
			},
			{
				Config:            testAccPolicyPackAttachmentConfig(),
				ResourceName:      "turbot_policy_pack_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccCheckAkasContain asserts the named akas list is non-empty and contains want, without
// depending on the order the API returns akas in.
func testAccCheckAkasContain(resourceName, attr, want string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		count := rs.Primary.Attributes[attr+".#"]
		if count == "" || count == "0" {
			return fmt.Errorf("%s.%s is empty; suppressIfAkaMatches would have nothing to match", resourceName, attr)
		}
		n, err := strconv.Atoi(count)
		if err != nil {
			return fmt.Errorf("could not parse %s.%s.#: %s", resourceName, attr, err)
		}
		for i := 0; i < n; i++ {
			if rs.Primary.Attributes[fmt.Sprintf("%s.%d", attr, i)] == want {
				return nil
			}
		}
		return fmt.Errorf("%s.%s does not contain %q", resourceName, attr, want)
	}
}

// TestAccPolicyPackAttachment_TargetDeleted covers the drift that Exists must survive: the
// attachment TARGET is removed out-of-band. Reading the attachment from the resource side means
// a missing target yields Not Found, and helper/schema aborts the whole refresh on any error from
// Exists — so this must map to "not attached" and drop the attachment from state instead.
//
// The target has to be deleted OUT-OF-BAND to reach that path. Simply removing it from the config
// does not: Terraform destroys the attachment before the folder, so the target still exists at
// every point Exists runs. Step 2's PreConfig therefore deletes the folder through the API client,
// behind Terraform's back, before the refresh that begins the step.
func TestAccPolicyPackAttachment_TargetDeleted(t *testing.T) {
	var targetId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyPackAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyPackAttachmentTargetDeletedConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyPackAttachmentExists("turbot_policy_pack_attachment.doomed"),
					// capture the target id so the next step can delete it out-of-band
					func(state *terraform.State) error {
						rs, ok := state.RootModule().Resources["turbot_folder.doomed"]
						if !ok {
							return fmt.Errorf("Not found: turbot_folder.doomed")
						}
						targetId = rs.Primary.ID
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*apiClient.Client)
					if err := client.DeleteResource(targetId); err != nil {
						t.Fatalf("could not delete target %s out-of-band: %s", targetId, err)
					}
				},
				// The attachment is still in state but its target is gone. The refresh that opens
				// this step calls Exists against the deleted target: it must report "not attached"
				// and drop the attachment, not error out.
				Config: testAccPolicyPackAttachmentTargetDeletedTeardownConfig(),
			},
		},
	})
}

func testAccPolicyPackAttachmentTargetDeletedConfig() string {
	return `
resource "turbot_policy_pack" "doomed" {
  filter      = "resourceType:181381985925765 $.turbot.tags.a:b"
  description = "Policy Pack target-deletion testing"
  title       = "policy_pack_target_deleted"
}

resource "turbot_folder" "doomed" {
  parent      = "tmod:@turbot/turbot#/"
  title       = "provider_test_target_deleted"
  description = "target that gets deleted out from under the attachment"
}

resource "turbot_policy_pack_attachment" "doomed" {
  policy_pack = turbot_policy_pack.doomed.id
  resource    = turbot_folder.doomed.id
}
`
}

func testAccPolicyPackAttachmentTargetDeletedTeardownConfig() string {
	return `
resource "turbot_policy_pack" "doomed" {
  filter      = "resourceType:181381985925765 $.turbot.tags.a:b"
  description = "Policy Pack target-deletion testing"
  title       = "policy_pack_target_deleted"
}
`
}
