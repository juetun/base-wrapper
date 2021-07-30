package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/juetun/base-wrapper/lib/base"
)

// Aes 加密解密
type Aes struct {
	Context *base.Context
}

// NewAes
func NewAes() (res *Aes) {
	res = &Aes{}
	return
}
func (r *Aes) Encryption(text string, aesKey string) (res string, err error) {
	var encrypted []byte
	encrypted, err = r.aesEncrypt([]byte(text), []byte(aesKey))
	res = base64.StdEncoding.EncodeToString(encrypted)
	return
}

func (r *Aes) Decrypt(encrypted, aesKey string) (res string, err error) {
	var origin []byte
	origin, err = r.aesDecrypt([]byte(encrypted), []byte(aesKey))
	res = string(origin)
	return
}

// AES加密
func (r *Aes) aesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = r.pKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypt := make([]byte, len(origData))
	blockMode.CryptBlocks(crypt, origData)
	return crypt, nil
}

// AES解密
func (r *Aes) aesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = r.pKCS7UnPadding(origData)
	return origData, nil
}
func (r *Aes) pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func (r *Aes) pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
