package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
	"github.com/turbot/terraform-provider-turbot/helpers"
)

// properties which must be passed to a create/update call
var watchProperties = []interface{}{"resource", "action", "filters", "favorite"}

func getWatchProperties() []interface{} {
	excludedProperties := []string{"resource"}
	return helpers.RemoveProperties(watchProperties, excludedProperties)
}

func resourceTurbotWatch() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotWatchCreate,
		Read:   resourceTurbotWatchRead,
		Update: resourceTurbotWatchUpdate,
		Delete: resourceTurbotWatchDelete,
		Exists: resourceTurbotWatchExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotWatchImport,
		},
		Schema: map[string]*schema.Schema{
			"resource": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filters": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"favorite": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"handler": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceTurbotWatchExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.WatchExists(id)
}

func resourceTurbotWatchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	// build map of watch properties
	input := mapFromResourceData(d, watchProperties)

	watch, err := client.CreateWatch(input)
	if err != nil {
		return err
	}

	// assign the id
	d.SetId(watch.Turbot.Id)

	return resourceTurbotWatchRead(d, meta)
}

func resourceTurbotWatchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	// build map of watch properties
	input := mapFromResourceData(d, getWatchProperties())
	input["id"] = id

	_, err := client.UpdateWatch(input)
	if err != nil {
		return err
	}

	return resourceTurbotWatchRead(d, meta)
}

func resourceTurbotWatchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	watch, err := client.ReadWatch(id)
	if err != nil {
		if errors.NotFoundError(err) {
			// watch was not found - clear id
			d.SetId("")
		}
		return err
	}

	// NOTE currently turbot accepts array of filters but only uses the first
	if len(watch.Filters) > 0 {
		d.Set("filters", watch.Filters)
	}
	d.Set("action", watch.Handler["action"])
	d.Set("favorite", watch.Turbot.FavoriteId)
	d.Set("handler", watch.Handler)
	d.Set("description", watch.Description)

	// If the resource is imported, set the resourceId
	// else set the value provided in the config
	_, propertySet := d.GetOk("resource")
	if !propertySet {
		d.Set("resource", watch.Turbot.ResourceId)
	}

	return nil
}

func resourceTurbotWatchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()
	err := client.DeleteWatch(id)
	if err != nil {
		return err
	}

	// clear the id to show we have deleted
	d.SetId("")

	return nil
}

func resourceTurbotWatchImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotWatchRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
