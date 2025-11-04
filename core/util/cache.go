package util

import (
	"encoding/json"
)

func Serialize(obj any) ([]byte, error) {
	return json.Marshal(obj)
}

func Deserialize(obj []byte, dest any) error {
	return json.Unmarshal(obj, dest)
}
