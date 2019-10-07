package turbot

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/encryption"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/iancoleman/strcase"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"sort"
)

var resourceMetadataProperties = []interface{}{"tags"}

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
			"body": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressIfBodyMatches,
			},
			"tags": {
				Type:     schema.TypeMap,
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
	parent := d.Get("parent").(string)
	resourceType := d.Get("type").(string)
	body := d.Get("body").(string)

	// populate turbot data
	mutationTurbotData := mapFromResourceData(d, resourceMetadataProperties)
	turbotMetadata, err := client.CreateResource(resourceType, parent, body, mutationTurbotData)
	if err != nil {
		return err
	}

	// set parent_akas property by loading resource and fetching the akas
	if err := storeAkas(turbotMetadata.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(turbotMetadata.Id)
	return nil
}

func resourceTurbotResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	// build required properties from body
	properties, err := propertiesFromBody(d.Get("body").(string))
	if err != nil {
		return fmt.Errorf("error retrieving properties from resource body: %s", err.Error())
	}

	resource, err := client.ReadResource(id, properties)
	if err != nil {
		if apiClient.NotFoundError(err) {
			// resource was not found - clear id
			d.SetId("")
		}
		return err
	}

	// rebuild body from the resource
	body, err := bodyFromProperties(resource.Data)
	if err != nil {
		return fmt.Errorf("error building resource body: %s", err.Error())
	}

	// assign results back into ResourceData

	// set parent_akas property by loading resource and fetching the akas
	if err := storeAkas(resource.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	d.Set("parent", resource.Turbot.ParentId)
	d.Set("body", body)
	return nil
}

func resourceTurbotResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	body := d.Get("body").(string)
	parent := d.Get("parent").(string)
	resourceType := d.Get("type").(string)
	id := d.Id()
	mutationTurbotData := mapFromResourceData(d, resourceMetadataProperties)
	turbotMetadata, err := client.UpdateResource(id, resourceType, parent, body, mutationTurbotData)
	if err != nil {
		return err
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

func suppressIfClientSecret(k, old, new string, d *schema.ResourceData) bool {
	return old != ""
}

// body is a json string
// apply standard formatting to old and new bodys then compare
func suppressIfBodyMatches(k, old, new string, d *schema.ResourceData) bool {
	if old == "" || new == "" {
		return false
	}

	oldFormatted := formatBody(old)
	newFormatted := formatBody(new)
	return oldFormatted == newFormatted
}

// given a json string, unmarshal into a map and return a map of alias ->  propertyName
func propertiesFromBody(body string) (map[string]string, error) {
	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return nil, err
	}
	var properties = map[string]string{}
	for k := range data {
		properties[k] = k
	}
	return properties, nil
}

// given a map of resource properties, marshal into a json string
func bodyFromProperties(d map[string]interface{}) (string, error) {
	body, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// apply standard formatting to the json body to enable easy diffing
func formatBody(body string) string {
	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		// ignore error and just return original body
		return body
	}
	body, err := bodyFromProperties(data)
	if err != nil {
		// ignore error and just return original body
		return body
	}
	return body

}

// construct a map of property values to pass to a graphql mutation
func mapFromResourceData(d *schema.ResourceData, properties []interface{}) map[string]interface{} {
	// each element in the 'properties' array is either a map defining explicit name mappings, or a string containing the terraform property name.
	// this is converted to the turbot property name by performing a snake case -> lowerCamelCase conversion
	// to build the output map:
	// 1) extract the value from ResourceData using the terraform property name
	// 2) add the property to a map using the turbot property name
	var propertyMap = map[string]interface{}{}
	for _, element := range properties {
		terraformToTurbotMap, ok := element.(map[string]string)
		// if terraformProperty is a map, perform explicit mapping and merge result with existing map
		if ok {
			mergeMaps(propertyMap, mapFromResourceDataWithPropertyMap(d, terraformToTurbotMap))
		} else {
			// otherwise perform automatic mapping from snake case (Terraform format) to lowerCamelCase (Turbot format).
			terraformProperty := element.(string)
			value, propertySet := d.GetOk(terraformProperty)
			// if property is set, map it
			if propertySet {
				var turbotProperty = strcase.ToLowerCamel(terraformProperty)
				propertyMap[turbotProperty] = value
			}
		}
	}
	return propertyMap
}

func mapFromResourceDataWithPropertyMap(d *schema.ResourceData, terraformToTurbotMap map[string]string) map[string]interface{} {
	var resourcePropertyMap = map[string]interface{}{}
	for terraform, turbot := range terraformToTurbotMap {
		// get schema for property
		value, propertySet := d.GetOk(terraform)
		if propertySet {
			resourcePropertyMap[turbot] = value
		}
	}
	return resourcePropertyMap
}

func mergeMaps(m1, m2 map[string]interface{}) {
	for k, v := range m2 {
		m1[k] = v
	}
}

// given a list of properties or property maps, remove the excluded properties
func removeProperties(properties []interface{}, excluded []string) []interface{} {
	for _, excludedProperty := range excluded {
		for i, element := range properties {
			// each element may be either a map, or a single property name
			terraformToTurbotMap, ok := element.(map[string]string)
			if ok {
				// if the element is a map, remove excluded items from map
				properties[i] = removePropertiesFromMap(terraformToTurbotMap, excluded)
			} else {
				// otherwise check if this property is excluded and remove if so
				if element.(string) == excludedProperty {
					properties = append(properties[:i], properties[i+1:]...)
					break
				}
			}
		}
	}
	return properties
}

// given a property list, remove the excluded properties
func removePropertiesFromMap(propertyMap map[string]string, excluded []string) map[string]string {
	var result = map[string]string{}
	for k, v := range propertyMap {
		if !sliceContains(excluded, k) {
			result[k] = v
		}
	}
	return result
}

// no native contains in golang :/
func sliceContains(s []string, searchTerm string) bool {
	i := sort.SearchStrings(s, searchTerm)
	return i < len(s) && s[i] == searchTerm

}

func encryptValue(pgpKey, value string) (string, string, error) {
	encryptionKey, err := encryption.RetrieveGPGKey(pgpKey)
	if err != nil {
		return "", "", err
	}
	fingerprint, encrypted, err := encryption.EncryptValue(encryptionKey, value, "Secret Key")
	if err != nil {
		return "", "", err
	}
	return fingerprint, encrypted, nil
}

func storeAkas(aka, propertyName string, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	akas, err := client.GetResourceAkas(aka)
	if err != nil {
		return err
	}
	// assign  akas
	d.Set(propertyName, akas)
	return nil
}
