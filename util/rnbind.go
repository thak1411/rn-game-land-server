package util

import (
	"encoding/json"
	"io"
)

func Bind(r io.ReadCloser, obj interface{}) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}
