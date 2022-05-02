package cipher

import "testing"

func TestAes(t *testing.T) {
	key := []byte("1234567890123456")
	plaintext := []byte("1234567890123456")
	ciphertext, err := AesEncrypt(plaintext, key)
	if err != nil {
		t.Error(err)
	}
	plaintext2, err := AesDecrypt(ciphertext, key)
	if err != nil {
		t.Error(err)
	}
	if string(plaintext) != string(plaintext2) {
		t.Error("plaintext != plaintext2")
	}
}
