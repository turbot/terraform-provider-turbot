package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"github.com/terraform-providers/terraform-provider-turbot/helpers"
)

// properties which must be passed to a create/update call
var smartFolderProperties = []interface{}{"title", "description", "parent", "filter"}

func getSmartFolderUpdateProperties() []interface{} {
	excludedProperties := []string{"parent"}
	return helpers.RemoveProperties(smartFolderProperties, excludedProperties)
}

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
				DiffSuppressFunc: suppressIfAkaMatches("parent_akas"),
			},
			//when doing a read, fetch the parent akas to use in suppressIfAkaMatches
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
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotSmartFolderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// build map of folder properties
	input := mapFromResourceData(d, smartFolderProperties)

	smartFolder, err := client.CreateSmartFolder(input)
	if err != nil {
		return err
	}

	// assign the id
	d.SetId(smartFolder.Turbot.Id)
	// TODO Remove Read call once schema changes are In.
	return resourceTurbotSmartFolderRead(d, meta)
}

func resourceTurbotSmartFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	// build map of folder properties
	input := mapFromResourceData(d, getSmartFolderUpdateProperties())
	input["id"] = id

	_, err := client.UpdateSmartFolder(input)
	if err != nil {
		return err
	}
	// set 'Read' Properties
	// TODO Remove Read call once schema changes are In.
	return resourceTurbotSmartFolderRead(d, meta)
}

func resourceTurbotSmartFolderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	smartFolder, err := client.ReadSmartFolder(id)
	if err != nil {
		if apiClient.NotFoundError(err) {
			// folder was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData
	// set parent_akas property by loading resource and fetching the akas
	if err := storeAkas(smartFolder.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// NOTE currently turbot accepts array of filters but only uses the first
	if len(smartFolder.Filters) > 0 {
		d.Set("filter", smartFolder.Filters[0])
	}
	d.Set("parent", smartFolder.Parent)
	d.Set("title", smartFolder.Title)
	d.Set("description", smartFolder.Description)

	return nil
}

func resourceTurbotSmartFolderDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotSmartFolderImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotSmartFolderRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
