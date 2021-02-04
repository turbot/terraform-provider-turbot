package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
)

// properties which must be passed to a create/update call
var folderDataProperties = []interface{}{"title", "description"}
var folderInputProperties = []interface{}{"parent", "tags", "akas"}

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
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressIfDescriptionNotPresent,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"akas": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTurbotFolderExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotFolderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	// build mutation input
	input := mapFromResourceData(d, folderInputProperties)
	input["data"] = mapFromResourceData(d, folderDataProperties)

	folder, err := client.CreateFolder(input)
	if err != nil {
		return err
	}

	// set parent_akas property by loading resource and fetching the akas
	if err := storeAkas(folder.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}

	// assign the id
	d.SetId(folder.Turbot.Id)
	// set FolderProperties the way we get in Read query
	d.Set("parent", folder.Parent)
	d.Set("title", folder.Title)
	d.Set("description", folder.Description)
	return nil
}

func resourceTurbotFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	// build mutation payload
	input := mapFromResourceData(d, folderInputProperties)
	input["data"] = mapFromResourceData(d, folderDataProperties)
	input["id"] = d.Id()

	folder, err := client.UpdateFolder(input)
	if err != nil {
		return err
	}
	// set FolderProperties the way we get in Read query
	d.Set("parent", folder.Parent)
	d.Set("title", folder.Title)
	d.Set("description", folder.Description)
	// set parent_akas property by loading resource and fetching the akas
	return storeAkas(folder.Turbot.ParentId, "parent_akas", d, meta)
}

func resourceTurbotFolderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	folder, err := client.ReadFolder(id)
	if err != nil {
		if errors.NotFoundError(err) {
			// folder was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData
	d.Set("parent", folder.Parent)
	d.Set("title", folder.Title)
	d.Set("description", folder.Description)
	d.Set("tags", folder.Turbot.Tags)
	d.Set("akas", folder.Turbot.Akas)
	// set parent_akas property by loading resource and fetching the akas
	return storeAkas(folder.Turbot.ParentId, "parent_akas", d, meta)
}

func resourceTurbotFolderDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotFolderImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotFolderRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

// folder resource type doesn't support deletion of description attribute
func suppressIfDescriptionNotPresent(_, old, new string, d *schema.ResourceData) bool {
	// if description is not in config, but in state-file
	if old != "" && new == "" {
		return true
	}
	return false
}
