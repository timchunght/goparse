package test

import (
	"encoding/json"
	// "testing"
)

func MapFromJSON(data []byte) map[string]interface{} {
	var result interface{}
	json.Unmarshal(data, &result)
	return result.(map[string]interface{})
}
