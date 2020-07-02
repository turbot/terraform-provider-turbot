package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"github.com/terraform-providers/terraform-provider-turbot/helpers"
)

var resourceProperties = []interface{}{"parent", "type", "tags", "akas"}

func getResourceUpdateProperties() []interface{} {
	excludedProperties := []string{"type"}
	return helpers.RemoveProperties(resourceProperties, excludedProperties)
}

func resourceTurbotResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotResourceCreate,
		Read:   resourceTurbotResourceRead,
		Update: resourceTurbotResourceUpdate,
		Delete: resourceTurbotResourceDelete,
		Exists: resourceTurbotResourceExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotResourceImport,
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
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressIfDataMatches,
			},
			"metadata": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressIfDataMatches,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func resourceTurbotResourceExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotResourceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	typeUri := d.Get("type")
	var err error

	// build input map to pass to mutation
	input, err := buildResourceInput(d, resourceProperties)
	if err != nil {
		return err
	}

	turbotMetadata, err := client.CreateResource(input)
	if err != nil {
		return err
	}

	// set parent_akas property by loading resource and fetching the akas
	if err := storeAkas(turbotMetadata.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(turbotMetadata.Id)
	// save the formatted data: this is to ensure the acceptance tests behave in a consistent way regardless of the ordering of the json data
	d.Set("data", helpers.FormatJson(d.Get("data").(string)))
	if metadata, ok := d.GetOk("metadata"); ok {
		d.Set("metadata", helpers.FormatJson(metadata.(string)))
	}
	d.Set("type", typeUri)
	return nil
}

func resourceTurbotResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	// build required properties from data.
	// properties is a map of property name -> property path

	var properties map[string]string = nil

	if _, ok := d.GetOk("data"); ok {
		var err error = nil
		properties, err = helpers.PropertyMapFromJson(d.Get("data").(string))
		if err != nil {
			return fmt.Errorf("error retrieving properties from resource data: %s", err.Error())
		}
	}

	resource, err := client.ReadResource(id, properties)
	if err != nil {
		if apiClient.NotFoundError(err) {
			// resource was not found - clear id
			d.SetId("")
		}
		return err
	}

	// rebuild data from the resource
	data, err := helpers.MapToJsonString(resource.Data)
	if err != nil {
		return fmt.Errorf("error building resource data: %s", err.Error())
	}

	// assign results back into ResourceData

	// set parent_akas property by loading resource and fetching the akas
	if err := storeAkas(resource.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}

	d.Set("parent", resource.Turbot.ParentId)
	d.Set("type", resource.Type.Uri)
	d.Set("tags", resource.Turbot.Tags)
	d.Set("data", data)
	return nil
}

func resourceTurbotResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// build input map to pass to mutation
	id := d.Id()
	input, err := buildResourceInput(d, getResourceUpdateProperties())
	if err != nil {
		return err
	}
	excludedPropertiesInUpdate, err := client.BuildPropertiesFromUpdateSchema(id, []interface{}{"updateSchema"})
	if err != nil {
		return err
	}
	input["data"], _ = buildDataUpdateProperties(d, excludedPropertiesInUpdate)
	input["id"] = d.Id()

	turbotMetadata, err := client.UpdateResource(input)
	if err != nil {
		return err
	}
	// save the formatted data: this is to ensure the acceptance tests behave in a consistent way regardless of the ordering of the json data
	d.Set("data", helpers.FormatJson(d.Get("data").(string)))
	if metadata, ok := d.GetOk("metadata"); ok {
		d.Set("metadata", helpers.FormatJson(metadata.(string)))
	}
	// set parent_akas property by loading resource and fetching the akas
	return storeAkas(turbotMetadata.ParentId, "parent_akas", d, meta)
}

func resourceTurbotResourceDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotResourceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotResourceRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func buildDataUpdateProperties(d *schema.ResourceData, properties []interface{}) (map[string]interface{}, error) {
	var err error
	// convert data from json string to map
	var dataMap map[string]interface{}
	dataString := d.Get("data").(string)
	if dataMap, err = helpers.JsonStringToMap(dataString); err != nil {
		return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal data: \n%s\nerror: %s", dataString, err.Error())
	}
	for _, element := range properties {
		if _, ok := dataMap[element.(string)]; ok {
			delete(dataMap, element.(string))
		}
	}

	return dataMap, nil
}

func buildResourceInput(d *schema.ResourceData, properties []interface{}) (map[string]interface{}, error) {
	var err error
	input := mapFromResourceData(d, properties)
	// convert data from json string to map
	dataString := d.Get("data").(string)
	if input["data"], err = helpers.JsonStringToMap(dataString); err != nil {
		return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal data: \n%s\nerror: %s", dataString, err.Error())
	}
	// convert metadata from json string to map (if present)
	if metadata, ok := d.GetOk("metadata"); ok {
		metadataString := metadata.(string)
		if input["metadata"], err = helpers.JsonStringToMap(metadataString); err != nil {
			return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal metadata: \n%s\nerror: %s", metadataString, err.Error())
		}
	}
	return input, nil
}

// the property in the config is an aka - however the state file will have an id.
// to perform a diff we also store the list of akas in state file
// if the new value of th eproperty exists in the akas list, then suppress diff
func suppressIfAkaMatches(propertyName string) func(k, old, new string, d *schema.ResourceData) bool {
	return func(k, old, new string, d *schema.ResourceData) bool {
		akasProperty, akasSet := d.GetOk(propertyName)
		// if parent_id has not been set yet, do not suppress the diff
		if !akasSet {
			return false
		}

		akas, ok := akasProperty.([]interface{})
		if !ok {
			return false
		}
		// if parentAkas contains 'new', suppress diff
		for _, aka := range akas {
			if aka.(string) == new {
				return true
			}
		}
		return false
	}

}

// data is a json string
// apply standard formatting to old and new data then compare
func suppressIfDataMatches(k, old, new string, d *schema.ResourceData) bool {
	if old == "" || new == "" {
		return false
	}

	oldFormatted := helpers.FormatJson(old)
	newFormatted := helpers.FormatJson(new)
	return oldFormatted == newFormatted
}
