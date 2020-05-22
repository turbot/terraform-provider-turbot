package helpers

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"reflect"
)

func ParseYamlString(value interface{}) (interface{}, error) {
	var temp interface{}

	s := InterfaceToString(value)

	err := yaml.Unmarshal([]byte(s), &temp)
	if err != nil {
		return s, err
	}

	if temp == "" {
		return value, nil
	}

	return temp, nil
}

// convert value to a string. If it is a complex type (object/array) then the diff calculation will use the value_source so the precise format is not critical
func InterfaceToString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

func InterfaceToYaml(value interface{}) (string, error) {
	if value == nil {
		return "", nil
	}
	if res, ok := value.(string); ok {
		return res, nil
	}

	data, err := yaml.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func YamlValuesAreEquivalent(yaml1, yaml2 string) (bool, error) {
	var yaml1intermediate, yaml2intermediate interface{}

	if err := yaml.Unmarshal([]byte(yaml1), &yaml1intermediate); err != nil {
		return false, fmt.Errorf("Error unmarshaling yaml string: %s", err)
	}

	if err := yaml.Unmarshal([]byte(yaml2), &yaml2intermediate); err != nil {
		return false, fmt.Errorf("Error unmarshaling yaml string: %s", err)
	}
	if reflect.DeepEqual(yaml1intermediate, yaml2intermediate) {
		return true, nil
	}
	return false, nil
}
