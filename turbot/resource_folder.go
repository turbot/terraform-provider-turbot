package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
)

func resourceTurbotFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotFolderCreate,
		Read:   resourceTurbotFolderRead,
		Update: resourceTurbotFolderCreate,
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
				// when doing a diff, the state file will contain the id of the parent bu tthe config contains the aka,
				// so we need custom diff code
				DiffSuppressFunc: supressIfParentAkaMatches,
			},
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
	title := d.Get("title").(string)
	description := d.Get("description").(string)

	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.CreateFolder(parentAka, title, description)
	if err != nil {
		return err
	}
	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(turbotMetadata.Id)

	return nil
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

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(d, meta); err != nil {
		return err
	}
	d.Set("parent", folder.Parent)
	d.Set("title", folder.Title)
	d.Set("description", folder.Description)

	return nil
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
