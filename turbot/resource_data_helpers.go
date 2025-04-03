package turbot

import (
	"reflect"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/iancoleman/strcase"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/helpers"
)

// given the resource data and a list of properties, construct a map of property values
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
			helpers.MergeMaps(propertyMap, mapFromResourceDataWithPropertyMap(d, terraformToTurbotMap))
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

// given the resource data and a map of properties (terraform property name -> output property name), construct a map of property values
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

// given a resource aka, fetch all akas for the resource and store in resourceData using 'propertyName'
func storeAkas(aka, propertyName string, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	akas, err := client.GetResourceAkas(aka)
	if err != nil {
		return err
	}
	// assign akas
	d.Set(propertyName, akas)
	return nil
}

// given an attribute, check for if the attribute contains a nil value
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
