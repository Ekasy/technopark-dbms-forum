package models

import "encoding/json"

func ToBytes(body interface{}) []byte {
	buf, err := json.Marshal(body)
	if err != nil {
		return []byte(`{"message": "could not make correct answer"}`)
	}
	return buf
}

type Error struct {
	Message string `json:"message"`
}
