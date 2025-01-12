package model

import (
	"encoding/json"
	"log"
)

func ToJSON(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to serialize data:", err)
		return ""
	}
	return string(bytes)
}

func FromJSONIntSlice(data string) []int {
	var result []int
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Println("Failed to deserialize data:", err)
	}
	return result
}

func FromJSONStringSlice(data string) []string {
	var result []string
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Println("Failed to deserialize data:", err)
	}
	return result
}
