package turbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestEncoding(t *testing.T) {
	var requestBody bytes.Buffer
	var queryObj = map[string]interface{}{
		"body": "test",
	}
	var variables = map[string]interface{}{
		"meta": "cool",
	}
	requestBodyObj := struct {
		Query     map[string]interface{} `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query:     queryObj,
		Variables: variables,
	}
	if err := json.NewEncoder(&requestBody).Encode(requestBodyObj); err != nil {
		fmt.Errorf("error: %s", err.Error())
	}
	var op interface{}
	json.NewDecoder(&requestBody).Decode(&op)
	log.Println("Decoded", requestBodyObj)
	log.Println("Decoded", op)
}
