package util

import (
	"math/rand"
	"time"
)

var randSeed = rand.New(rand.NewSource(time.Now().UnixNano()))

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func AddPS(x string) string {
	return "THIS_IS_RN'S_PREFIX_" + x
}

func GenGuestName() string {
	res := make([]byte, 6)
	for i := range res {
		res[i] = charSet[randSeed.Intn(len(charSet))]
	}
	return "GUEST_" + string(res)
}
