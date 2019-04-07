package util

import (
	"crypto/hmac"
	"crypto/sha256"
)

func MakeHmac(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}
