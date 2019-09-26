package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
)

// properties which must be passed to a create/update call
// TODO add filters here once we are consistent with the db
var smartFolderProperties = []string{"title", "description"}

func resourceTurbotSmartFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotSmartFolderCreate,
		Read:   resourceTurbotSmartFolderRead,
		Update: resourceTurbotSmartFolderUpdate,
		Delete: resourceTurbotSmartFolderDelete,
		Exists: resourceTurbotSmartFolderExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotSmartFolderImport,
		},
		Schema: map[string]*schema.Schema{
			//aka of the parent resource
			"parent": {
				Type:     schema.TypeString,
				Required: true,
				// when doing a diff, the state file will contain the id of the parent but the config contains the aka,
				// so we need custom diff code
				DiffSuppressFunc: suppressIfParentAkaMatches,
			},
			//when doing a read, fetch the parent akas to use in suppressIfParentAkaMatches()
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
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceTurbotSmartFolderExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiclient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotSmartFolderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	// build map of folder properties
	data := mapFromResourceData(d, smartFolderProperties)
	// TODO currently turbot accepts array of filters but only uses the first
	if filter := d.Get("filter").(string); filter != "" {
		data["filters"] = []string{filter}
	}

	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.CreateSmartFolder(parentAka, data)
	if err != nil {
		return err
	}

	// set parent_akas property by loading resource resource and fetching the akas
	parentAkas, err := client.GetResourceAkas(turbotMetadata.ParentId)
	if err != nil {
		return err
	}
	// assign parent_akas
	d.Set("parent_akas", parentAkas)

	// assign the id
	d.SetId(turbotMetadata.Id)
	return nil
}

func resourceTurbotSmartFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	id := d.Id()

	// build map of folder properties
	data := mapFromResourceData(d, smartFolderProperties)
	// TODO currently turbot accepts array of filters but only uses the first
	if filter := d.Get("filter").(string); filter != "" {
		data["filters"] = []string{filter}
	}
	//data["filters"] = []string{d.Get("filter").(string)}

	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.UpdateSmartFolder(id, parentAka, data)
	if err != nil {
		return err
	}
	// set parent_akas property by loading resource resource and fetching the akas
	parentAkas, err := client.GetResourceAkas(turbotMetadata.ParentId)
	if err != nil {
		return err
	}
	// assign parent_akas
	d.Set("parent_akas", parentAkas)
	return nil
}

func resourceTurbotSmartFolderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	id := d.Id()

	smartFolder, err := client.ReadSmartFolder(id)
	if err != nil {
		if apiclient.NotFoundError(err) {
			// folder was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData
	// set parent_akas property by loading resource resource and fetching the akas
	parentAkas, err := client.GetResourceAkas(smartFolder.Turbot.ParentId)
	if err != nil {
		return err
	}
	if len(smartFolder.Filters) > 0 {
		d.Set("filter", smartFolder.Filters[0])
	}
	// assign parent_akas
	d.Set("parent_akas", parentAkas)
	d.Set("parent_id", smartFolder.Parent)
	// TODO currently turbot accepts array of filters but only uses the first

	d.Set("title", smartFolder.Title)
	d.Set("description", smartFolder.Description)

	return nil
}

func resourceTurbotSmartFolderDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotSmartFolderImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotSmartFolderRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
