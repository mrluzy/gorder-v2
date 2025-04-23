package util

import "encoding/json"

func MarshalString(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	return string(bytes), err
}
