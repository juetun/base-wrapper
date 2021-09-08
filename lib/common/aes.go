package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/wumansgy/goEncrypt"
)

// Aes 加密解密
type Aes struct {
	Context *base.Context
}

// NewAes 加密操作函数
func NewAes() (res *Aes) {
	res = &Aes{}
	return
}

// EncryptionCtr 加密操作
func (r *Aes) EncryptionCtr(text string, aesKey string) (res string, err error) {
	var cryptText []byte
	if cryptText, err = goEncrypt.AesCtrEncrypt([]byte(text), []byte(aesKey), nil); err != nil {
		return
	}
	res = base64.StdEncoding.EncodeToString(cryptText)
	return
}

func (r *Aes) DecryptCtr(text, aesKey string) (res string, err error) {
	var encrypted, resBt []byte
	if encrypted, err = base64.StdEncoding.DecodeString(text); err != nil {
		return
	}
	// 传入密文和自己定义的密钥，需要和加密的密钥一样，不一样会报错 可以自己传入初始化向量,如果不传就使用默认的初始化向量,16字节
	if resBt, err = goEncrypt.AesCtrDecrypt(encrypted, []byte(aesKey), nil); err != nil {
		return
	}
	res = string(resBt)
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
