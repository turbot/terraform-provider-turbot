package helpers

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Describes a tag.
type Tag struct {
	_     struct{} `type:"structure"`
	Key   *string  `locationName:"key" type:"string"`
	Value *string  `locationName:"value" type:"string"`
}

// tagsSchema returns the schema to use for tags.
//
func TagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}

// tagsFromMap returns the tags for the given map of data.
func TagsFromMap(m map[string]string) []*Tag {
	result := make([]*Tag, 0, len(m))
	for k, v := range m {
		t := &Tag{
			Key:   &k,
			Value: &v,
		}
		result = append(result, t)
	}

	return result
}

// tagsToMap turns the list of tags into a map.
func TagsToMap(ts []*Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range ts {
		result[*t.Key] = *t.Value
	}

	return result
}
