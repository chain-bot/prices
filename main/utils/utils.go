package utils

import (
	"encoding/json"
	"log"
)

func PrettyJSON(obj interface{}) string {
	json, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(json)
}
