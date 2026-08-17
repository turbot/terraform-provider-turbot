package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"strings"
)

func resourceTurbotSmartFolderAttachemnt() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotSmartFolderAttachmentCreate,
		Read:   resourceTurbotSmartFolderAttachmentRead,
		Delete: resourceTurbotSmartFolderAttachmentDelete,
		Exists: resourceTurbotSmartFolderAttachmentExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotSmartFolderAttachmentImport,
		},
		Schema: map[string]*schema.Schema{
			"resource": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressIfAkaMatches("resource_akas"),
			},
			"smart_folder": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressIfAkaMatches("smart_folder_akas"),
			},
			"resource_akas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Stores the smart folder's AKAs so suppressIfAkaMatches can suppress diffs
			// when the user provides an AKA but the state holds the resolved numeric ID.
			"smart_folder_akas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTurbotSmartFolderAttachmentExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	smartFolderId, resource := parseSmartFolderId(d.Id())

	// Check the attachment from the RESOURCE side rather than reading the smart folder and
	// enumerating everything attached to it. Reading the smart folder requires a grant on it
	// (smart folders live at the Turbot root by default), whereas the caller necessarily
	// holds permissions on the attachment target. See apiClient/policy_pack.go.
	attached, err := client.PolicyPackAttached(resource, smartFolderId)
	if err != nil {
		// A deleted target takes its attachments with it. helper/schema aborts the whole refresh on
		// any error from Exists, so returning one here would force the operator to
		// `terraform state rm`; returning false instead drops the attachment and lets it replan.
		//
		// Decided from the response shape via apiClient.ErrTargetNotFound, NOT from matching "not
		// found" in the error text - an unrelated not-found would otherwise drop a live attachment
		// out of state.
		if apiClient.IsTargetNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("error reading smart folder attachment: %s", err.Error())
	}
	return attached, nil
}

func resourceTurbotSmartFolderAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	resource := d.Get("resource").(string)
	smartFolder := d.Get("smart_folder").(string)

	// Resolve the smart_folder AKA or ID to its numeric Turbot ID.
	// The attachSmartFolders mutation requires numeric IDs — AKA strings cause "not eligible for attachment" errors.
	// ReadPolicyPackIdentity uses `policyPack(id:)`, which accepts either form and, unlike the
	// generic `resource(id:)`, does not require a grant on the smart folder itself. A smart folder
	// and a policy pack are the same underlying resource type. See apiClient/policy_pack.go.
	sfIdentity, err := client.ReadPolicyPackIdentity(smartFolder, "smart folder")
	if err != nil {
		return err
	}
	resolvedSmartFolderId := sfIdentity.Id
	if resolvedSmartFolderId == "" {
		return fmt.Errorf("smart folder %q resolved to an empty ID", smartFolder)
	}

	input := map[string]interface{}{
		"resource":     resource,
		"smartFolders": resolvedSmartFolderId,
	}

	_, err = client.CreateSmartFolderAttachment(input)
	if err != nil {
		return err
	}

	// Store resource AKAs for DiffSuppressFunc on the resource field
	if err := storeAkas(resource, "resource_akas", d, meta); err != nil {
		return err
	}
	// Reuse AKAs from the already-fetched smart folder identity to avoid a second round-trip
	smartFolderAkas := sfIdentity.Akas
	if len(smartFolderAkas) == 0 {
		smartFolderAkas = []string{resolvedSmartFolderId}
	}
	d.Set("smart_folder_akas", smartFolderAkas)

	// Always store the resolved numeric ID in state and the state ID so parseSmartFolderId
	// (which splits on the first underscore) works correctly for all input formats.
	d.SetId(buildId(resolvedSmartFolderId, resource))
	d.Set("resource", resource)
	d.Set("smart_folder", resolvedSmartFolderId)
	return nil
}

func resourceTurbotSmartFolderAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// NOTE: This will not be called if the attachment does not exist
	smartFolder, resource := parseSmartFolderId(d.Id())

	turbotResource, err := client.ReadResource(resource, nil)
	if err != nil {
		return err
	}
	// set resource_akas property by loading resource and fetching the akas
	if err := storeAkas(turbotResource.Turbot.Id, "resource_akas", d, meta); err != nil {
		return err
	}
	// set smart_folder_akas property for DiffSuppressFunc. Read via `policyPack(id:)` rather than
	// storeAkas (which goes through `resource(id:)`) so this does not require a grant on the smart
	// folder — see apiClient/policy_pack.go.
	sfIdentity, err := client.ReadPolicyPackIdentity(smartFolder, "smart folder")
	if err != nil {
		return err
	}
	smartFolderAkas := sfIdentity.Akas
	if len(smartFolderAkas) == 0 {
		smartFolderAkas = []string{smartFolder}
	}
	d.Set("smart_folder_akas", smartFolderAkas)
	// assign results directly back into ResourceData
	d.Set("resource", resource)
	d.Set("smart_folder", smartFolder)
	return nil
}

func resourceTurbotSmartFolderAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	smartFolder, resource := parseSmartFolderId(d.Id())
	input := map[string]interface{}{
		"resource":     resource,
		"smartFolders": smartFolder,
	}
	err := client.DeleteSmartFolderAttachment(input)
	if err != nil {
		return err
	}

	// clear the id to show we have deleted
	d.SetId("")
	return nil
}

func resourceTurbotSmartFolderAttachmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotSmartFolderAttachmentRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func buildId(smartFolder, resource string) string {
	return smartFolder + "_" + resource
}

func parseSmartFolderId(id string) (smartFolder, resource string) {
	segments := strings.Split(id, "_")
	smartFolder = segments[0]
	resource = strings.Join(segments[1:], "_")
	return
}
