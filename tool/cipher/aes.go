package cipher
import (
	"fmt"
	"crypto/cipher"
	"crypto/aes"
	"bytes"
	"encoding/base64"
	)
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
	return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	cryptData := make([]byte, len(origData))
	blockMode.CryptBlocks(cryptData, origData)
	return cryptData, nil
}
func AesDecrypt(cryptData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
	return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cryptData))
	blockMode.CryptBlocks(origData, cryptData)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}