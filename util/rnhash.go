package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func Encrypt(data, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(data + salt))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}
