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

// func BindClaims(iClaims, obj interface{}) error {
// 	if iClaims == nil {
// 		return errors.New("nil pointer can't binding")
// 	}
// 	json.Unmarshal([]byte(iClaims), obj)
// }
