package helpers

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestRemoveProperties(t *testing.T) {
	type test struct {
		name       string
		properties []interface{}
		excluded   []string
		expected   []interface{}
	}
	tests := []test{
		test{
			"No exclusions",
			[]interface{}{"a", "b", "c"},
			[]string{},
			[]interface{}{"a", "b", "c"},
		},
		test{
			"String exclusions",
			[]interface{}{"a", "b", "c"},
			[]string{"a"},
			[]interface{}{"b", "c"},
		},
		test{
			"All excluded",
			[]interface{}{"a", "b", "c"},
			[]string{"a", "b", "c"},
			[]interface{}(nil),
		},
		test{
			"Map exclusion",
			[]interface{}{"a", "b", map[string]string{"c": "C", "d": "D"}},
			[]string{"c"},
			[]interface{}{"a", "b", map[string]string{"d": "D"}},
		},
		test{
			"2 map exclusions",
			[]interface{}{"a", "b", map[string]string{"c": "C", "d": "D"}, map[string]string{"e": "E", "f": "F"}},
			[]string{"c", "f"},
			[]interface{}{"a", "b", map[string]string{"d": "D"}, map[string]string{"e": "E"}},
		},
		test{
			"No matching exclusions",
			[]interface{}{"a", "b", "c"},
			[]string{"d"},
			[]interface{}{"a", "b", "c"},
		},
		test{
			"No matching exclusions with map",
			[]interface{}{"a", "b", map[string]string{"c": "C", "d": "D"}},
			[]string{"e"},
			[]interface{}{"a", "b", map[string]string{"c": "C", "d": "D"}},
		},
	}
	for _, test := range tests {
		log.Println(test.name)
		result := RemoveProperties(test.properties, test.excluded)
		assert.Equal(t, test.expected, result)
	}
}
