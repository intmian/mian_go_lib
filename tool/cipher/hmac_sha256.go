package cipher

import (
	"crypto/hmac"
	"crypto/sha256"
)

func HmacSha256Sign(secret string, data string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return h.Sum(nil)
}
