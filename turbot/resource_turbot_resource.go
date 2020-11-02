package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"github.com/terraform-providers/terraform-provider-turbot/errorHandler"
	"github.com/terraform-providers/terraform-provider-turbot/helpers"
	"strings"
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
		if errorHandler.NotFoundError(err) {
			// resource was not found - clear id
			d.SetId("")
		}
		return err
	}
	var data string
	if _, ok := d.GetOk("data"); ok {
		// if data is set, only include the properties that are specified in the resource config
		data, err := getStringValueForKey(d,"data",resource.Data)
		if err != nil {
			return err
		}
		d.Set("data", data)
	} else if _, ok := d.GetOk("full_data"); ok {
		// if full_data is set, include all data read from API
		if data, err = helpers.MapToJsonString(resource.Data); err != nil {
			return fmt.Errorf("error retrieving data properties: %s", err.Error())
		}
		d.Set("full_data", data)
	}

	// In the import case, we won't have this
	var metadata string
	if _, ok := d.GetOk("metadata"); ok {
		// if metadata is set, only include the properties that are specified in the resource config
		if metadata, err = getStringValueForKey(d,"metadata",resource.Turbot.Custom);err != nil {
			return fmt.Errorf("error retrieving metadata properties: %s", err.Error())
		}
		d.Set("metadata", metadata)
	} else if _, ok := d.GetOk("full_metadata"); ok {
		// if full_metadata is set, include all data read from API
		if metadata, err = helpers.MapToJsonString(resource.Turbot.Custom); err != nil {
			return fmt.Errorf("error retrieving metadata properties: %s", err.Error())
		}
		d.Set("full_metadata", metadata)
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
	var dataProperty string
	var ok bool
	if _, ok = d.GetOk("data"); ok {
		dataProperty = "data"
	} else if _, ok = d.GetOk("full_data"); ok {
		dataProperty = "full_data"
	}
	if ok {
		input["data"], err = buildUpdatePayloadForData(d, client, dataProperty)
		if err != nil {
			return err
		}
	}
	// Identify metadata property (metadata/full_netadata)
	var metaProperty string
	if _, ok = d.GetOk("metadata"); ok {
		metaProperty = "metadata"
	} else if _, ok = d.GetOk("full_metadata"); ok {
		metaProperty = "full_metadata"
	}
	if ok {
		input["metadata"], err = buildUpdatePayloadForMetadata(d, metaProperty)
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
	var ok bool
	var dataString, metaDataString interface{}
	input := mapFromResourceData(d, properties)
	if dataString, ok = d.GetOk("data"); !ok {
		dataString, ok = d.GetOk("full_data")
	}
	if ok {
		if input["data"], err = helpers.JsonStringToMap(dataString.(string)); err != nil {
			return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal data: \n%s\nerror: %s", dataString, err.Error())
		}
	}
	// convert metadata from json string to map (if present)
	if metaDataString, ok = d.GetOk("metadata"); !ok {
		metaDataString, ok = d.GetOk("full_metadata")
	}
	if ok {
		if input["metadata"], err = helpers.JsonStringToMap(metaDataString.(string)); err != nil {
			return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal data: \n%s\nerror: %s", dataString, err.Error())
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

// - build a map from the data or full_data property (specified by 'key' parameter)
// - add a `nil` value for deleted properties
// - remove any properties disallowed by the updateSchema
func buildUpdatePayloadForData(d *schema.ResourceData, client *apiClient.Client, key string) (map[string]interface{},error) {
	var err error
	dataMap, err := markPropertiesForDeletion(d,key)
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
	metaData, err := markPropertiesForDeletion(d,key)
	if err != nil {
		return nil ,err
	}
	return metaData, nil
}

func markPropertiesForDeletion(d *schema.ResourceData, key string) (map[string]interface{}, error) {
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
			// NOTE: for folder we cannot currently delete the description property
			shouldDelete := !(getResourceType(d.Get("type").(string)) == "folder" &&  key == "description")
			if _, ok := oldContent[key.(string)]; ok && shouldDelete {
				newContent[key.(string)] = nil
			}
		}
	}
	return newContent, nil
}
// - get properties for a given key from config
// - build a map only including the properties fetched from config
// - convert map to string
func getStringValueForKey(d *schema.ResourceData, key string, readResponse map[string]interface{}) (string,error) {
	propertiesOfKey, err := getPropertiesFromConfig(d, key)
	if err != nil {
		return "",err
	}
	metadata := buildResourceMapFromProperties(readResponse, propertiesOfKey)
	metadataString, err := helpers.MapToJsonString(metadata)
	if err != nil {
		return "",fmt.Errorf("error building resource data: %s", err.Error())
	}
	return metadataString, nil
}

func getResourceType(uri string) string {
	splitsOfUri := strings.Split(uri,"/")
	// underflow check
	var typeUri string
	if len(splitsOfUri) > 0 {
		typeUri = splitsOfUri[len(splitsOfUri)-1]
	}
	return typeUri
}