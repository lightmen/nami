package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func AesEncrypt(data []byte, key string) ([]byte, error) {
	bk := []byte(key)
	block, err := aes.NewCipher(bk)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCEncrypter(block, bk[:blockSize])

	padBytes := pkcs7Padding(data, blockSize)
	cryptData := make([]byte, len(padBytes))
	blockMode.CryptBlocks(cryptData, padBytes)

	return cryptData, nil
}

func AesDecrypt(data []byte, key string) ([]byte, error) {
	bk := []byte(key)
	block, err := aes.NewCipher(bk)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, bk[:blockSize])

	cryptData := make([]byte, len(data))
	blockMode.CryptBlocks(cryptData, data)
	cryptData, err = pkcs7UnPadding(cryptData)
	if err != nil {
		return nil, err
	}

	return cryptData, nil
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7UnPadding(data []byte) ([]byte, error) {
	ln := len(data)
	if ln == 0 {
		return nil, errors.New("decryption data error")
	}
	unPadding := int(data[ln-1])
	return data[:(ln - unPadding)], nil
}
