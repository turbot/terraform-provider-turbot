package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"github.com/terraform-providers/terraform-provider-turbot/helpers"
)

var fileProperties = []interface{}{"parent", "tags", "akas"}
var metadataProperties = []interface{}{"title", "description"}

func getFileUpdateProperties() []interface{} {
	excludedProperties := []string{"type"}
	return helpers.RemoveProperties(resourceProperties, excludedProperties)
}

func resourceTurbotFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotFileCreate,
		Read:   resourceTurbotFileRead,
		Update: resourceTurbotFileUpdate,
		Delete: resourceTurbotFileDelete,
		Exists: resourceTurbotFileExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotFileImport,
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
				Type:     schema.TypeString,
				Optional: true,
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

func resourceTurbotFileExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotFileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	description := d.Get("description")
	var err error
	// build input map to pass to mutation
	input, err := buildFileInput(d, fileProperties)
	// set type property
	input["type"] = "tmod:@turbot/turbot#/resource/types/file"
	// data should be object - strict
	// if data( object must ), override the title
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
	d.Set("title", d.Get("title"))
	d.Set("description", description)
	return nil
}

func resourceTurbotFileRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("title", resource.Turbot.Title)
	d.Set("description", resource.Turbot.Description)
	d.Set("type", resource.Type.Uri)
	d.Set("data", data)
	return nil
}

func resourceTurbotFileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// build input map to pass to mutation
	id := d.Id()
	input, err := buildFileInput(d, getFileUpdateProperties())
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
	d.Set("title", turbotMetadata.Title)
	d.Set("description", turbotMetadata.Description)
	// set parent_akas property by loading resource and fetching the akas
	return storeAkas(turbotMetadata.ParentId, "parent_akas", d, meta)
}

func resourceTurbotFileDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotFileImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotResourceRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func buildFileInput(d *schema.ResourceData, properties []interface{}) (map[string]interface{}, error) {
	var err error
	input := mapFromResourceData(d, properties)
	// convert data from json string to map
	dataString := d.Get("data").(string)
	if input["data"], err = helpers.JsonStringToMap(dataString); err != nil {
		return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal data: \n%s\nerror: %s", dataString, err.Error())
	}
	// convert metadata from json string to map (if present)
	// insert top level `title` and `description` property inside metadata
	if metadata, ok := d.GetOk("metadata"); ok {
		metadataString := metadata.(string)
		if input["metadata"], err = helpers.JsonStringToMap(metadataString); err != nil {
			return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal metadata: \n%s\nerror: %s", metadataString, err.Error())
		}
		metadataMap := input["metadata"].(map[string]interface{})
		for _, element := range metadataProperties {
			metadataProperty := element.(string)
			topLevelValue, propertySet := d.GetOk(metadataProperty)
			// if property is set, map it
			if propertySet {
				if value, ok := metadataMap[metadataProperty]; ok {
					if value != topLevelValue {
						return nil, fmt.Errorf("error title mismatch, failed to pass different title as top level: %s and metadata title:%s ", topLevelValue, value)
					}
				} else {
					metadataMap[metadataProperty] = topLevelValue
				}
			}
		}
		input["metadata"] = metadataMap
	}
	return input, nil
}
