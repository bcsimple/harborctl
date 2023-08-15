package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
)

func init() {
	key = []byte(secret)
	hash := md5.Sum(key)
	// aes key必须是16,24,32字节
	key = hash[:]
}

func Encrypt(rawData string) (string, error) {
	encryptData, err := CFBEncrypt(key, []byte(rawData))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptData), nil
}

func Decrypt(encryptedDataStr string) (string, error) {
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedDataStr)
	if err != nil {
		return "", err
	}
	rawData, err := CFBDecrypt(key, encryptedData)
	if err != nil {
		return "", err
	}
	return string(rawData), nil
}

func CFBEncrypt(key []byte, rawData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cfbEncrypter := cipher.NewCFBEncrypter(block, iv)
	result := make([]byte, len(rawData))
	cfbEncrypter.XORKeyStream(result, rawData)
	return result, nil
}

func CFBDecrypt(key []byte, encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cfbDecrypter := cipher.NewCFBDecrypter(block, iv)
	result := make([]byte, len(encryptedData))
	cfbDecrypter.XORKeyStream(result, encryptedData)
	return result, nil
}
