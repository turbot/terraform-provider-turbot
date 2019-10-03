package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
)

// properties which must be passed to a create/update call
var folderDataProperties = []interface{}{"title", "description"}
var folderMetadataProperties = []interface{}{"tags"}

func resourceTurbotFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotFolderCreate,
		Read:   resourceTurbotFolderRead,
		Update: resourceTurbotFolderUpdate,
		Delete: resourceTurbotFolderDelete,
		Exists: resourceTurbotFolderExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotFolderImport,
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
			"description": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceTurbotFolderExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiclient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotFolderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	// build mutation payload
	payload := map[string]map[string]interface{}{
		"data":       mapFromResourceData(d, folderDataProperties),
		"turbotData": mapFromResourceData(d, folderMetadataProperties),
	}

	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.CreateFolder(parentAka, payload)
	if err != nil {
		return err
	}

	// set parent_akas property by loading resource and fetching the akas
	if err := storeAkas(turbotMetadata.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}

	// assign the id
	d.SetId(turbotMetadata.Id)
	return nil
}

func resourceTurbotFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	id := d.Id()

	// build mutation payload
	payload := map[string]map[string]interface{}{
		"data":       mapFromResourceData(d, folderDataProperties),
		"turbotData": mapFromResourceData(d, folderMetadataProperties),
	}
	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.UpdateFolder(id, parentAka, payload)
	if err != nil {
		return err
	}
	// set parent_akas property by loading resource and fetching the akas
	return storeAkas(turbotMetadata.ParentId, "parent_akas", d, meta)
}

func resourceTurbotFolderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	id := d.Id()

	folder, err := client.ReadFolder(id)
	if err != nil {
		if apiclient.NotFoundError(err) {
			// folder was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData
	d.Set("parent", folder.Parent)
	d.Set("title", folder.Title)
	d.Set("description", folder.Description)
	// set parent_akas property by loading resource and fetching the akas
	return storeAkas(folder.Turbot.ParentId, "parent_akas", d, meta)
}

func resourceTurbotFolderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	id := d.Id()
	err := client.DeleteResource(id)
	if err != nil {
		return err
	}

	// clear the id to show we have deleted
	d.SetId("")

	return nil
}

func resourceTurbotFolderImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotFolderRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
