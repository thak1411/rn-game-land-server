package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func encrypt(data, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(data + salt))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No Args")
		return
	}
	data := os.Args[1]

	fmt.Println("INPUT:", data)
	fmt.Println(encrypt(data, "admin_salt"))
}
