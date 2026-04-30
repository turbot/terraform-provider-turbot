package turbot

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
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
	// execute api call
	smartFolder, err := client.ReadSmartFolder(policyPackId)
	if err != nil {
		return false, fmt.Errorf("error reading policy pack: %s", err.Error())
	}

	// find resource aka in list of attached resources
	for _, attachedResource := range smartFolder.AttachedResources.Items {
		if resource == attachedResource.Turbot.Id {
			return true, nil
		}
		for _, aka := range attachedResource.Turbot.Akas {
			if aka == resource {
				return true, nil
			}
		}
	}
	return false, nil
}

func resourceTurbotPolicyPackAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	resource := d.Get("resource").(string)
	policyPack := d.Get("policy_pack").(string)

	// Resolve the policy_pack AKA or ID to its numeric Turbot ID.
	// The attachSmartFolders mutation requires numeric IDs — AKA strings cause "not eligible for attachment" errors.
	// ReadResource uses a typed turbot { id } query which correctly populates Turbot.Id,
	// unlike ReadSmartFolder which uses get(path:"turbot") and does not parse Id reliably.
	policyPackResource, err := client.ReadResource(policyPack, nil)
	if err != nil {
		return fmt.Errorf("error reading policy pack %q: %s", policyPack, err.Error())
	}
	resolvedPolicyPackId := policyPackResource.Turbot.Id

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
	// Store policy pack AKAs for DiffSuppressFunc on the policy_pack field
	if err := storeAkas(resolvedPolicyPackId, "policy_pack_akas", d, meta); err != nil {
		return err
	}

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
	// set policy_pack_akas property for DiffSuppressFunc
	if err := storeAkas(policyPack, "policy_pack_akas", d, meta); err != nil {
		return err
	}
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
