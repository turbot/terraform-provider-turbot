package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
)

// these are the properties which must be passed to a create/update call
var localDirectoryProperties = []interface{}{"title", "profile_id_template", "description"}
var localDirectoryMetadataProperties = []interface{}{"tags"}

func resourceTurbotLocalDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotLocalDirectoryCreate,
		Read:   resourceTurbotLocalDirectoryRead,
		Update: resourceTurbotLocalDirectoryUpdate,
		Delete: resourceTurbotLocalDirectoryDelete,
		Exists: resourceTurbotLocalDirectoryExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotLocalDirectoryImport,
		},
		Schema: map[string]*schema.Schema{
			// aka of the parent resource
			"parent": {
				Type:     schema.TypeString,
				Required: true,
				// when doing a diff, the state file will contain the id of the parent but the config contains the aka,
				// so we need custom diff code
				DiffSuppressFunc: suppressIfAkaMatches("parent_akas"),
			},
			// when doing a read, fetch the parent akas to use in suppressIfAkaMatches
			"parent_akas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"profile_id_template": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"directory_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTurbotLocalDirectoryExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotLocalDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	parentAka := d.Get("parent").(string)
	// build mutation payload
	payload := map[string]map[string]interface{}{
		"data":       helpers.MapFromResourceData(d, localDirectoryProperties),
		"turbotData": helpers.MapFromResourceData(d, localDirectoryMetadataProperties),
	}
	// set computed properties
	payload["data"]["status"] = "Active"
	payload["data"]["directoryType"] = "local"
	turbotMetadata, err := client.CreateLocalDirectory(parentAka, payload)
	if err != nil {
		return err
	}

	// set parent_akas property by loading resource and fetching the akas
	if err := helpers.StoreAkas(turbotMetadata.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(turbotMetadata.Id)
	d.Set("status", payload["data"]["status"])
	d.Set("directory_type", payload["data"]["directoryType"])
	return nil
}

func resourceTurbotLocalDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	parentAka := d.Get("parent").(string)
	id := d.Id()

	// build mutation payload
	payload := map[string]map[string]interface{}{
		"data":       helpers.MapFromResourceData(d, localDirectoryProperties),
		"turbotData": helpers.MapFromResourceData(d, localDirectoryMetadataProperties),
	}
	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.UpdateLocalDirectory(id, parentAka, payload)
	if err != nil {
		return err
	}
	// set parent_akas property by loading resource and fetching the akas
	return helpers.StoreAkas(turbotMetadata.ParentId, "parent_akas", d, meta)
}

func resourceTurbotLocalDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	localDirectory, err := client.ReadLocalDirectory(id)
	if err != nil {
		if apiClient.NotFoundError(err) {
			// local directoery was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData
	d.Set("parent", localDirectory.Parent)
	d.Set("title", localDirectory.Title)
	d.Set("status", localDirectory.Status)
	d.Set("directory_type", localDirectory.DirectoryType)
	// set parent_akas property by loading resource and fetching the akas
	return helpers.StoreAkas(localDirectory.Turbot.ParentId, "parent_akas", d, meta)
}

func resourceTurbotLocalDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()
	err := client.DeleteResource(id)
	if err != nil {
		return err
	}

	// clear the id to show we have deleted
	d.SetId("")
	return nil
}

func resourceTurbotLocalDirectoryImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotLocalDirectoryRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
