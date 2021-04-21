package TechCustomers

import (
	"encoding/json"
)

type BodyText struct {
	body []string `json:"BodyText"`
}

func CustomerList(jsonData []byte) {
	var jsonStruct BodyText

	json.Unmarshal([]byte(jsonData), &jsonStruct)
}
