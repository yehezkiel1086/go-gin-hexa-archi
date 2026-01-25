package util

import (
	"encoding/json"
	"fmt"
)

func GenerateCacheKey(prefix string, param any) string {
	return fmt.Sprintf("%s:%v", prefix, param)
}

func GenerateCacheKeyParams(params ...any) string {
	var str string

	for i, param := range params {
		str += fmt.Sprintf("%v", param)

		last := len(params) - 1
		if i != last {
			str += "-"
		}
	}

	return str
}

func Serialize(data any) ([]byte, error) {
	return json.Marshal(data)
}

func Deserialize(data []byte, res any) error {
	return json.Unmarshal(data, res)
}
