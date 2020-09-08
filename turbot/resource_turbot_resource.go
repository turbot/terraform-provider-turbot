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
				Optional:         true,
				DiffSuppressFunc: suppressIfDataMatches,
			},
			"metadata": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressIfDataMatches,
			},
			"full_data": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressIfDataMatches,
				ConflictsWith:    []string{"data"},
			},
			"full_metadata": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressIfDataMatches,
				ConflictsWith:    []string{"metadata"},
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
	if data, ok := d.GetOk("data"); ok {
		d.Set("data", helpers.FormatJson(data.(string)))
	}
	if metadata, ok := d.GetOk("metadata"); ok {
		d.Set("metadata", helpers.FormatJson(metadata.(string)))
	}
	if fullData, ok := d.GetOk("full_data"); ok {
		d.Set("full_data", helpers.FormatJson(fullData.(string)))
	}
	if fullMetadata, ok := d.GetOk("full_metadata"); ok {
		d.Set("full_metadata", helpers.FormatJson(fullMetadata.(string)))
	}
	d.Set("type", typeUri)
	return nil
}

func resourceTurbotResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()
	var err error
	// read full resource
	resource, err := client.ReadFullResource(id)
	if err != nil {
		if apiClient.NotFoundError(err) {
			// resource was not found - clear id
			d.SetId("")
		}
		return err
	}
	// build required properties from schema for data and metadata.
	// properties is a map of property name -> property path
	var dataProperties map[string]string
	var metadataProperties map[string]string

	// include the properties that are specified in the resource config
	if _, ok := d.GetOk("data"); ok {
		dataProperties, err = getPropertiesFromConfig(d, "data")
		if err != nil {
			return err
		}
		data := buildResourceMapFromProperties(resource.Data, dataProperties)
		dataString, err := helpers.MapToJsonString(data)
		if err != nil {
			return fmt.Errorf("error building resource data: %s", err.Error())
		}
		d.Set("data", dataString)
	}
	if _, ok := d.GetOk("metadata"); ok {
		metadataProperties, err = getPropertiesFromConfig(d, "metadata")
		if err != nil {
			return err
		}
		metadata := buildResourceMapFromProperties(resource.Turbot.Custom, metadataProperties)
		metadataString, err := helpers.MapToJsonString(metadata)
		if err != nil {
			return fmt.Errorf("error building resource data: %s", err.Error())
		}
		d.Set("metadata", metadataString)
	}

	// full_data and full_metadata attributes are directly read from readFullResource response
	// we don't exclude any keys from the read response.
	if _, ok := d.GetOk("full_metadata"); ok {
		if metadata, err := helpers.MapToJsonString(resource.Turbot.Custom); err != nil {
			d.Set("full_metadata", metadata)
		}
	}
	if _, ok := d.GetOk("full_data"); ok {
		if data, err := helpers.MapToJsonString(resource.Data); err != nil {
			d.Set("full_data", data)
		}
	}

	// set parent_akas property by loading resource and fetching the akas
	if err := storeAkas(resource.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	d.Set("parent", resource.Turbot.ParentId)
	d.Set("type", resource.Type.Uri)
	d.Set("tags", resource.Turbot.Tags)
	return nil
}

func resourceTurbotResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// build input map to pass to mutation
	id := d.Id()
	// build mutation input data by parsing the resource schema and
	// excluding top level properties which must not be sent to update - `type`
	input, err := buildResourceInput(d, getResourceUpdateProperties())
	if err != nil {
		return err
	}
	// Identify data property (data/full_data)
	if _, ok := d.GetOk("data"); ok {
		// remove the keys that were previously there but not in new config
		input["data"], err = buildUpdatePayloadForData(d, client, "data")
		if err != nil {
			return err
		}
	} else if _, ok := d.GetOk("full_data"); ok {
		// remove the keys that were previously there but not in the new config
		// also set the keys to nil that were set externally
		// by other means
		input["data"], err = buildUpdatePayloadForData(d, client, "full_data")
		if err != nil {
			return err
		}
	}
	// Identify metadata property (metadata/full_netadata)
	if _, ok := d.GetOk("metadata"); ok {
		// remove the keys that were previously there but not in new config
		input["metadata"], err = buildUpdatePayloadForMetadata(d, "metadata")
		if err != nil {
			return err
		}
	} else if _, ok := d.GetOk("full_data"); ok {
		// remove the keys that were previously there but not in the new config
		// also set the keys to nil that were set externally,  by other means
		input["metadata"], err = buildUpdatePayloadForMetadata(d, "full_metadata")
		if err != nil {
			return err
		}
	}
	input["id"] = id
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

func buildDataUpdateProperties(data map[string]interface{}, forbiddenProperties []interface{}) map[string]interface{} {
	for _, element := range forbiddenProperties {
		if _, ok := data[element.(string)]; ok {
			delete(data, element.(string))
		}
	}
	return data
}

func buildResourceInput(d *schema.ResourceData, properties []interface{}) (map[string]interface{}, error) {
	var err error
	input := mapFromResourceData(d, properties)
	// convert data from json string to map
	if dataString, ok := d.GetOk("data"); ok {
		if input["data"], err = helpers.JsonStringToMap(dataString.(string)); err != nil {
			return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal data: \n%s\nerror: %s", dataString, err.Error())
		}
	} else
		if dataString, ok := d.GetOk("full_data"); ok {
			if input["data"], err = helpers.JsonStringToMap(dataString.(string)); err != nil {
			return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal data: \n%s\nerror: %s", dataString, err.Error())
		}
	}
	// convert metadata from json string to map (if present)
	if metadata, ok := d.GetOk("metadata"); ok {
		if input["metadata"], err = helpers.JsonStringToMap(metadata.(string)); err != nil {
			return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal metadata: \n%s\nerror: %s", metadata, err.Error())
		}
	} else
		if metadata, ok := d.GetOk("full_metadata"); ok {
			if input["metadata"], err = helpers.JsonStringToMap(metadata.(string)); err != nil {
			return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal metadata: \n%s\nerror: %s", metadata, err.Error())
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

func getPropertiesFromConfig(d *schema.ResourceData, key string) (map[string]string, error) {
	var properties map[string]string = nil
	var err error = nil
	if keyValue, ok := d.GetOk(key); ok {
		if properties, err = helpers.PropertyMapFromJson(keyValue.(string)); err!= nil {
			return nil, fmt.Errorf("error retrieving properties: %s", err.Error())
		}
	}
	return properties, nil
}

func buildResourceMapFromProperties(input map[string]interface{}, properties map[string]string) map[string]interface{} {
	for key, _ := range input {
		// delete external keys from response data
		if _, ok := properties[key]; !ok {
			delete(input, key)
		}
	}
	return input
}

// buildUpdatePayload(propertyName) -
// - build a map from the data or full_data property (specified by 'key' parameter)
// - add a `nil` value for deleted properties
// - remove any properties disallowed by the updateSchema
func buildUpdatePayloadForData(d *schema.ResourceData, client *apiClient.Client, key string) (map[string]interface{},error) {
	var err error
	dataMap, err := setOldPropertiesToNull(d,key)
	if err != nil {
		return nil ,err
	}
	excludedPropertiesInUpdate, err := client.BuildPropertiesFromUpdateSchema(d.Id(), []interface{}{"updateSchema"})
	if err != nil {
		return nil ,err
	}
	// build data object by excluding forbidden properties from either `data` and `full_data`
	dataMap = buildDataUpdateProperties(dataMap, excludedPropertiesInUpdate)
	return dataMap, nil
}

func buildUpdatePayloadForMetadata(d *schema.ResourceData, key string) (map[string]interface{},error)  {
	metaData, err := setOldPropertiesToNull(d,key)
	if err != nil {
		return nil ,err
	}
	return metaData, nil
}

func setOldPropertiesToNull(d *schema.ResourceData, key string) (map[string]interface{}, error) {
	var oldContent, newContent map[string]interface{}
	var err error
	// fetch old(state-file) and new(config) content
	if old, new := d.GetChange(key); old != nil {
		if oldContent, err = helpers.JsonStringToMap(old.(string)); err != nil {
			return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal content: \n%s\nerror: %s", old.(string), err.Error())
		}
		if newContent, err = helpers.JsonStringToMap(new.(string)); err != nil {
			return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal content: \n%s\nerror: %s", new.(string), err.Error())
		}
		// extract keys from old content not in new
		excludeContentProperties := helpers.GetOldMapProperties(oldContent, newContent)
		for _, key := range excludeContentProperties {
			// set keys of old content to `nil` in new content
			// any property which doesn't exist in config is set to nil
			if _, ok := oldContent[key.(string)]; ok {
				newContent[key.(string)] = nil
			}
		}
	}
	return newContent, nil
}