package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"log"
	"strings"
)

var localDirectoryProperties = []string{"title", "profile_id_template"}

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
		},
	}
}

func resourceTurbotLocalDirectoryExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiclient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotLocalDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	// build map of local directory properties
	data := mapFromResourceData(d, localDirectoryProperties)
	data["status"] = "New"
	data["directoryType"] = "local"
	turbotMetadata, err := client.CreateLocalDirectory(parentAka, data)
	if err != nil {
		return err
	}

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(turbotMetadata.ParentId, d, meta); err != nil {
		return err
	}

	// assign the id
	d.SetId(turbotMetadata.Id)
	d.Set("status", data["status"])
	d.Set("directoryType", data["directoryType"])
	return nil
}

func existingLocalDirectoryError(existingLocalDirectories []apiclient.LocalDirectory, title, parentAka string) error {
	// build array of existing folder ids
	var ids []string
	for _, f := range existingLocalDirectories {
		if f.Turbot.Id != "" {
			ids = append(ids, f.Turbot.Id)
		}
	}
	// TODO extract terraform name
	var directoryString, idString string

	if len(ids) > 1 {
		directoryString = "folders"
		idString = "ids"
	} else {
		directoryString = "a folder"
		idString = "id"
	}
	return fmt.Errorf("Cannot create Local Directory '%s' with parent '%s' as %s of that name already exists in that location, with %s: %s. To manage an existing Turbot local directory using Terraform, import it using command 'terraform import <resource_address> <id>'",
		title, parentAka, directoryString, idString, strings.Join(ids, ","))
}

func resourceTurbotLocalDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	id := d.Id()

	// build map of local directory properties
	data := mapFromResourceData(d, folderProperties)
	log.Println("[INFO} resourceTurbotLocalDirectoryUpdate", id, parentAka)
	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.UpdateDirectory(id, parentAka, data)
	if err != nil {
		return err
	}
	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(turbotMetadata.ParentId, d, meta); err != nil {
		return err
	}
	return nil
}

func resourceTurbotLocalDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	id := d.Id()

	localDirectory, err := client.ReadLocalDirectory(id)
	if err != nil {
		if apiclient.NotFoundError(err) {
			// local directoery was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(localDirectory.Turbot.ParentId, d, meta); err != nil {
		return err
	}
	d.Set("parent", localDirectory.Parent)
	d.Set("title", localDirectory.Title)
	return nil
}

func resourceTurbotLocalDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotLocalDirectoryImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotLocalDirectoryRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
