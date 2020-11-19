package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"github.com/terraform-providers/terraform-provider-turbot/errors"
	"github.com/terraform-providers/terraform-provider-turbot/helpers"
	"strings"
)

// these are the properties which must be passed to a create/update call
var turbotDirectoryInputProperties = []interface{}{"parent", "status", "title", "description", "profile_id_template", "server", "tags"}

func getTurbotDirectoryUpdateProperties() []interface{} {
	excludedProperties := []string{"profile_id_template", "server"}
	return helpers.RemoveProperties(turbotDirectoryInputProperties, excludedProperties)
}
func resourceTurbotTurbotDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotTurbotDirectoryCreate,
		Read:   resourceTurbotTurbotDirectoryRead,
		Update: resourceTurbotTurbotDirectoryUpdate,
		Delete: resourceTurbotTurbotDirectoryDelete,
		Exists: resourceTurbotTurbotDirectoryExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotTurbotDirectoryImport,
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
			"server": {
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
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceTurbotTurbotDirectoryExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotTurbotDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// build mutation input
	input := mapFromResourceData(d, turbotDirectoryInputProperties)
	// set computed properties
	input["status"] = "ACTIVE"

	// do create
	turbotDirectory, err := client.CreateTurbotDirectory(input)
	if err != nil {
		return err
	}

	// set parent_akas property by loading resource and fetching the akas
	if err := storeAkas(turbotDirectory.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(turbotDirectory.Turbot.Id)
	// assign properties coming back from create graphQl API
	d.Set("parent", turbotDirectory.Turbot.ParentId)
	d.Set("title", turbotDirectory.Title)
	// Set the values from Resource Data
	d.Set("status", input["status"])
	return nil
}

func resourceTurbotTurbotDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	turbotDirectory, err := client.ReadTurbotDirectory(id)
	if err != nil {
		if errors.NotFoundError(err) {
			// local directoery was not found - clear id
			d.SetId("")
		}
		return err
	}
	// assign results back into ResourceData

	d.Set("title", turbotDirectory.Title)
	d.Set("description", turbotDirectory.Description)
	d.Set("status", strings.ToUpper(turbotDirectory.Status))
	d.Set("parent", turbotDirectory.Turbot.ParentId)
	d.Set("profile_id_template", turbotDirectory.ProfileIdTemplate)
	d.Set("tags", turbotDirectory.Turbot.Tags)
	d.Set("server", turbotDirectory.Server)
	// set parent_akas property by loading resource and fetching the akas
	return storeAkas(turbotDirectory.Turbot.ParentId, "parent_akas", d, meta)
}

func resourceTurbotTurbotDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// build mutation payload
	input := mapFromResourceData(d, getTurbotDirectoryUpdateProperties())
	input["id"] = d.Id()

	// do update
	turbotDirectory, err := client.UpdateTurbotDirectory(input)
	if err != nil {
		return err
	}

	// assign properties coming back from update graphQl API
	d.Set("parent", turbotDirectory.Turbot.ParentId)
	d.Set("title", turbotDirectory.Title)
	d.Set("status", turbotDirectory.Status)
	// set parent_akas property by loading resource and fetching the akas
	return storeAkas(turbotDirectory.Turbot.ParentId, "parent_akas", d, meta)
}

func resourceTurbotTurbotDirectoryImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotTurbotDirectoryRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func resourceTurbotTurbotDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
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
