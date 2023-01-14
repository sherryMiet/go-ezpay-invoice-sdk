package ezpay_invoice

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

func AesCBCEncrypt(rawData, key []byte, iv []byte) (string, error) {
	rawData = PKCS7Padding(rawData)
	ciphertext := make([]byte, len(rawData))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, rawData)

	return hex.EncodeToString(ciphertext), nil
}
func AesCBCDecrypt(rawData, key []byte, iv []byte) ([]byte, error) {

	ciphertext := make([]byte, len(rawData))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, rawData)

	ciphertext = PKCS7UnPadding(ciphertext)
	return ciphertext, nil
}
func PKCS7Padding(ciphertext []byte) []byte {
	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}
