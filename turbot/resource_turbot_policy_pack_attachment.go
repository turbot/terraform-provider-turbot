package turbot

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
)

func resourceTurbotPolicyPackAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotPolicyPackAttachmentCreate,
		Read:   resourceTurbotPolicyPackAttachmentRead,
		Delete: resourceTurbotPolicyPackAttachmentDelete,
		Exists: resourceTurbotPolicyPackAttachmentExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotPolicyPackAttachmentImport,
		},
		Schema: map[string]*schema.Schema{
			"resource": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressIfAkaMatches("resource_akas"),
			},
			"policy_pack": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressIfAkaMatches("policy_pack_akas"),
			},
			"resource_akas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Stores the policy pack's AKAs so suppressIfAkaMatches can suppress diffs
			// when the user provides an AKA but the state holds the resolved numeric ID.
			"policy_pack_akas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTurbotPolicyPackAttachmentExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	policyPackId, resource := parsePolicyPackId(d.Id())

	// Check the attachment from the RESOURCE side rather than reading the policy pack and
	// enumerating everything attached to it. Reading the pack requires a grant on the pack
	// (packs live at the Turbot root by default), whereas the caller necessarily holds permissions on
	// the attachment target. This also keeps the answer scoped to this one resource instead
	// of depending on a pack-wide list that may span the whole hierarchy.
	attached, err := client.PolicyPackAttached(resource, policyPackId)
	if err != nil {
		// A deleted target takes its attachments with it. helper/schema aborts the whole refresh
		// on any error from Exists, so returning one here would force the operator to
		// `terraform state rm`; returning false instead drops the attachment and lets it replan.
		if errors.NotFoundError(err) {
			return false, nil
		}
		return false, fmt.Errorf("error reading policy pack attachment: %s", err.Error())
	}
	return attached, nil
}

func resourceTurbotPolicyPackAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	resource := d.Get("resource").(string)
	policyPack := d.Get("policy_pack").(string)

	// Resolve the policy_pack AKA or ID to its numeric Turbot ID.
	// The attachSmartFolders mutation requires numeric IDs — AKA strings cause "not eligible for attachment" errors.
	// ReadPolicyPackIdentity uses `policyPack(id:)`, which accepts either form and, unlike the
	// generic `resource(id:)`, does not require a grant on the pack itself. See apiClient/policy_pack.go.
	policyPackIdentity, err := client.ReadPolicyPackIdentity(policyPack, "policy pack")
	if err != nil {
		return err
	}
	resolvedPolicyPackId := policyPackIdentity.Id
	if resolvedPolicyPackId == "" {
		return fmt.Errorf("policy pack %q resolved to an empty ID", policyPack)
	}

	input := map[string]interface{}{
		"resource":     resource,
		"smartFolders": resolvedPolicyPackId,
	}

	_, err = client.CreateSmartFolderAttachment(input)
	if err != nil {
		return err
	}

	// Store resource AKAs for DiffSuppressFunc on the resource field
	if err := storeAkas(resource, "resource_akas", d, meta); err != nil {
		return err
	}
	// Reuse AKAs from the already-fetched policy pack identity to avoid a second round-trip
	policyPackAkas := policyPackIdentity.Akas
	if len(policyPackAkas) == 0 {
		policyPackAkas = []string{resolvedPolicyPackId}
	}
	d.Set("policy_pack_akas", policyPackAkas)

	// Always store the resolved numeric ID in state and the state ID so parsePolicyPackId
	// (which splits on the first underscore) works correctly for all input formats.
	d.SetId(buildPolicyPackId(resolvedPolicyPackId, resource))
	d.Set("resource", resource)
	d.Set("policy_pack", resolvedPolicyPackId)
	return nil
}

func resourceTurbotPolicyPackAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// NOTE: This will not be called if the attachment does not exist
	policyPack, resource := parsePolicyPackId(d.Id())

	turbotResource, err := client.ReadResource(resource, nil)
	if err != nil {
		return err
	}
	// set resource_akas property by loading resource and fetching the akas
	if err := storeAkas(turbotResource.Turbot.Id, "resource_akas", d, meta); err != nil {
		return err
	}
	// set policy_pack_akas property for DiffSuppressFunc. Read via `policyPack(id:)` rather
	// than storeAkas (which goes through `resource(id:)`) so this does not require a grant on
	// the pack — see apiClient/policy_pack.go.
	policyPackIdentity, err := client.ReadPolicyPackIdentity(policyPack, "policy pack")
	if err != nil {
		return err
	}
	policyPackAkas := policyPackIdentity.Akas
	if len(policyPackAkas) == 0 {
		policyPackAkas = []string{policyPack}
	}
	d.Set("policy_pack_akas", policyPackAkas)
	// assign results directly back into ResourceData
	d.Set("resource", resource)
	d.Set("policy_pack", policyPack)
	return nil
}

func resourceTurbotPolicyPackAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	policyPack, resource := parsePolicyPackId(d.Id())
	input := map[string]interface{}{
		"resource":     resource,
		"smartFolders": policyPack,
	}
	err := client.DeleteSmartFolderAttachment(input)
	if err != nil {
		return err
	}

	// clear the id to show we have deleted
	d.SetId("")
	return nil
}

func resourceTurbotPolicyPackAttachmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotPolicyPackAttachmentRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func buildPolicyPackId(policyPack, resource string) string {
	return policyPack + "_" + resource
}

func parsePolicyPackId(id string) (policyPack, resource string) {
	// Get the index of the first underscore
	index := strings.Index(id, "_")
	policyPack = id[:index]
	resource = id[index+1:]
	return
}
