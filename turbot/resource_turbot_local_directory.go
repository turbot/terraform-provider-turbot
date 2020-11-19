package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"github.com/terraform-providers/terraform-provider-turbot/errors"
	"github.com/terraform-providers/terraform-provider-turbot/helpers"
	"strings"
)

// input properties which must be passed to a create/update call
var localDirectoryInputProperties = []interface{}{"title", "profile_id_template", "parent", "description", "tags"}

// exclude properties from input map to make a update call
func getLocalDirectoryUpdateProperties() []interface{} {
	excludedProperties := []string{"profile_id_template", "tags"}
	return helpers.RemoveProperties(localDirectoryInputProperties, excludedProperties)
}

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

	// build mutation input

	input := mapFromResourceData(d, localDirectoryInputProperties)
	input["status"] = "ACTIVE"

	localDirectory, err := client.CreateLocalDirectory(input)
	if err != nil {
		return err
	}

	// set parent_akas property by loading resource and fetching the akas
	if err := storeAkas(localDirectory.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(localDirectory.Turbot.Id)
	// assign properties coming back from create graphQl API
	d.Set("parent", localDirectory.Parent)
	d.Set("title", localDirectory.Title)
	d.Set("status", strings.ToUpper(localDirectory.Status))
	d.Set("directory_type", localDirectory.DirectoryType)
	// Set the values from Resource Data
	return nil
}

func resourceTurbotLocalDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	localDirectory, err := client.ReadLocalDirectory(id)
	if err != nil {
		if errors.NotFoundError(err) {
			// local directory was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData
	d.Set("parent", localDirectory.Parent)
	d.Set("title", localDirectory.Title)
	d.Set("description", localDirectory.Description)
	d.Set("status", strings.ToUpper(localDirectory.Status))
	d.Set("profile_id_template", localDirectory.ProfileIdTemplate)
	d.Set("directory_type", localDirectory.DirectoryType)
	d.Set("tags", localDirectory.Turbot.Tags)
	// set parent_akas property by loading resource and fetching the akas
	return storeAkas(localDirectory.Turbot.ParentId, "parent_akas", d, meta)
}

func resourceTurbotLocalDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	// build mutation payload
	input := mapFromResourceData(d, getLocalDirectoryUpdateProperties())
	input["id"] = d.Id()
	// do update
	localDirectory, err := client.UpdateLocalDirectory(input)
	if err != nil {
		return err
	}

	// assign properties coming back from update graphQl API
	d.Set("parent", localDirectory.Parent)
	d.Set("title", localDirectory.Title)
	d.Set("status", strings.ToUpper(localDirectory.Status))
	d.Set("directory_type", localDirectory.DirectoryType)
	// set parent_akas property by loading resource and fetching the akas
	return storeAkas(localDirectory.Turbot.ParentId, "parent_akas", d, meta)
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
