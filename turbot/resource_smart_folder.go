package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
)

// properties which must be passed to a create/update call
var smartFolderProperties = map[string] interface{}{"title", "description", "rules"}

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
			// aka of the parent resource
			// "parent": {
			// 	Type:     schema.TypeString,
			// 	Required: true,
			// 	// when doing a diff, the state file will contain the id of the parent bu tthe config contains the aka,
			// 	// so we need custom diff code
			// 	DiffSuppressFunc: supressIfParentAkaMatches,
			// },
			// when doing a read, fetch the parent akas to use in supressIfParentAkaMatches()
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
			"rules":{
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			}

		},
	}
}

func resourceTurbotSmartFolderExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiclient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

// func resourceTurbotSmartFolderCreate(d *schema.ResourceData, meta interface{}) error {
// 	client := meta.(*apiclient.Client)
// 	parentAka := d.Get("parent").(string)
// 	// build map of folder properties
// 	data := mapFromResourceData(d, folderProperties)
// 	// create folder returns turbot resource metadata containing the id
// 	turbotMetadata, err := client.CreateFolder(parentAka, data)
// 	if err != nil {
// 		return err
// 	}

// 	// set parent_akas property by loading parent resource and fetching the akas
// 	if err = setParentAkas(turbotMetadata.ParentId, d, meta); err != nil {
// 		return err
// 	}

// 	// assign the id
// 	d.SetId(turbotMetadata.Id)

// 	return nil
// }

func resourceTurbotSmartFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	id := d.Id()

	// build map of folder properties
	data := mapFromResourceData(d, folderProperties)

	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.UpdateFolder(id, parentAka, data)
	if err != nil {
		return err
	}
	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(turbotMetadata.ParentId, d, meta); err != nil {
		return err
	}
	return nil
}

func resourceTurbotSmartFolderRead(d *schema.ResourceData, meta interface{}) error {
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

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(folder.Turbot.ParentId, d, meta); err != nil {
		return err
	}
	d.Set("rules", folder.Rules)
	d.Set("title", folder.Title)
	d.Set("description", folder.Description)

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
