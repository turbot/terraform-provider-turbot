package helpers

import (
	"encoding/json"
	"github.com/hashicorp/terraform/helper/encryption"
	"github.com/iancoleman/strcase"
	"reflect"
	"sort"
)

func MergeMaps(m1, m2 map[string]interface{}) {
	for k, v := range m2 {
		m1[k] = v
	}
}

// given a list of properties or property maps, remove the excluded properties
func RemoveProperties(properties []interface{}, excluded []string) []interface{} {
	for _, excludedProperty := range excluded {
		for i, element := range properties {
			// each element may be either a map, or a single property name
			terraformToTurbotMap, ok := element.(map[string]interface{})
			if ok {
				// if the element is a map, remove excluded items from map
				properties[i] = RemovePropertiesFromMap(terraformToTurbotMap, excluded)
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
func RemovePropertiesFromMap(propertyMap map[string]interface{}, excluded []string) map[string]interface{} {
	var result = map[string]interface{}{}
	for k, v := range propertyMap {
		if !SliceContains(excluded, k) {
			result[k] = v
		}
	}
	return result
}

// no native contains in golang :/
func SliceContains(s []string, searchTerm string) bool {
	i := sort.SearchStrings(s, searchTerm)
	return i < len(s) && s[i] == searchTerm

}

func EncryptValue(pgpKey, value string) (string, string, error) {
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

func ConvertToJsonString(data map[string]interface{}) (string, error) {
	dataBytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return "", err
	}
	jsonData := string(dataBytes)
	return jsonData, nil
}

// convert a map[string]interface{} to a map[string]string byt json encoding any non string fields
func ConvertToStringMap(data map[string]interface{}) (map[string]string, error) {
	var outputMap = map[string]string{}

	for k, v := range data {
		if v != nil {
			if reflect.TypeOf(v).String() != "string" {
				jsonBytes, err := json.MarshalIndent(v, "", " ")
				if err != nil {
					return nil, err
				}
				v = string(jsonBytes)
			}
			outputMap[k] = v.(string)
		}
	}
	return outputMap, nil
}

// convert map keys to snake case
func ConvertMapKeysToSnakeCase(data map[string]interface{}) map[string]interface{} {
	var outputMap = map[string]interface{}{}
	for k, v := range data {
		outputMap[strcase.ToLowerCamel(k)] = v
	}
	return outputMap
}
