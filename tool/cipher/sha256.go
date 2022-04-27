package cipher

import (
	"crypto/sha256"
	"fmt"
)

func Sha2562String(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}
