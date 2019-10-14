package helpers

import (
	"encoding/json"
	"github.com/hashicorp/terraform/helper/encryption"
	"reflect"
	"sort"
)

func MergeMaps(m1, m2 map[string]interface{}) {
	for k, v := range m2 {
		m1[k] = v
	}
}

// given a list of items which may each be either a property or property map, remove the excluded properties
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

// TODO update to use delete
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

func MapToJsonString(data map[string]interface{}) (string, error) {
	dataBytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return "", err
	}
	jsonData := string(dataBytes)
	return jsonData, nil
}

func JsonStringToMap(dataString string) (map[string]interface{}, error) {
	var data = make(map[string]interface{})
	if err := json.Unmarshal([]byte(dataString), &data); err != nil {
		return nil, err
	}
	return data, nil
}

// apply standard formatting to a json string by unmarshalling into a map then marshalling back to JSON
func FormatJson(body string) string {
	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		// ignore error and just return original body
		return body
	}
	body, err := MapToJsonString(data)
	if err != nil {
		// ignore error and just return original body
		return body
	}
	return body

}

// given a json representation of an object, build a map of the property names: property alias -> property path
func PropertyMapFromJson(body string) (map[string]string, error) {
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

// convert a map[string]interface{} to a map[string]string by json encoding any non string fields
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
