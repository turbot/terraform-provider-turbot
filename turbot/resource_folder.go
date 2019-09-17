package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"log"
	"strings"
)

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

	// first check if the folder exists - search by parent and foldere title
	existingFolders, err := client.FindFolder(title, parentAka)
	if err != nil {
		return err
	}
	if len(existingFolders) > 0 {
		return existingFolderError(existingFolders, title, parentAka)
	}

	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.CreateFolder(parentAka, title, description)
	if err != nil {
		return err
	}

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(turbotMetadata.ParentId, d, meta); err != nil {
		return err
	}

	// assign the id
	d.SetId(turbotMetadata.Id)

	return nil
}

func existingFolderError(existingFolders []apiclient.Folder, title, parentAka string) error {
	// build array of existing folder ids
	var ids []string
	for _, f := range existingFolders {
		if f.Turbot.Id != "" {
			ids = append(ids, f.Turbot.Id)
		}
	}
	// TODO extract terraform name
	var folderString, idString string

	if len(ids) > 1 {
		folderString = "folders"
		idString = "ids"
	} else {
		folderString = "a folder"
		idString = "id"
	}
	return fmt.Errorf("Cannot create folder '%s' with parent '%s' as %s of that name already exists in that location, with %s: %s. To manage an existing Turbot folder using Terraform, import it using command 'terraform import <resource_address> <id>'",
		title, parentAka, folderString, idString, strings.Join(ids, ","))
}

func resourceTurbotFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	title := d.Get("title").(string)
	description := d.Get("description").(string)
	id := d.Id()

	log.Println("[INFO} resourceTurbotFolderUpdate", id, parentAka)
	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.UpdateFolder(id, parentAka, title, description)
	if err != nil {
		return err
	}
	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(turbotMetadata.ParentId, d, meta); err != nil {
		return err
	}
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
	if err = setParentAkas(folder.Turbot.ParentId, d, meta); err != nil {
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
