package model

import (
	"encoding/json"
	"log"
)

func toJSON(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to serialize data:", err)
		return ""
	}
	return string(bytes)
}

func fromJSONIntSlice(data string) []int {
	var result []int
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Println("Failed to deserialize data:", err)
	}
	return result
}

func fromJSONStringSlice(data string) []string {
	var result []string
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Println("Failed to deserialize data:", err)
	}
	return result
}
