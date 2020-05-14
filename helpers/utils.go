package helpers

import (
	"fmt"
	"github.com/go-yaml/yaml"
)

func ParseYamlString(value interface{}) (interface{}, error) {
	var temp interface{}

	s := SettingValueToString(value)
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
func SettingValueToString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}
