package util

import (
	"encoding/json"
	"io"
)

func BindBody(body io.ReadCloser, obj interface{}) error {
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}

func BindJson(jsn []byte, obj interface{}) error {
	if err := json.Unmarshal(jsn, obj); err != nil {
		return err
	}
	return nil
}
