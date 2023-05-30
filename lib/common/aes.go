package common

import (
	"bytes"
	aes2 "crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/wumansgy/goEncrypt/aes"
	"strings"
)

const (
	ChatEncryptionKey = "wuMansgygoaesctr" // 秘钥长度为16的倍数
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
	var (
		block cipher.Block
		data  = []byte(text)
		key   = []byte(aesKey)
	)
	if block, err = aes2.NewCipher(key); err != nil {
		return
	}
	//判断加密快的大小
	blockSize := block.BlockSize()
	//填充
	encryptBytes := pkcs7Padding(data, blockSize)

	var cryptText []byte

	if cryptText, err = aes.AesCtrEncrypt(encryptBytes, key, []byte(ChatEncryptionKey)); err != nil {
		return
	}
	res = base64.StdEncoding.EncodeToString(cryptText)
	return
}

// pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("加密字符串错误！")
	}
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

func (r *Aes) DecryptCtr(text, aesKey string) (res string, err error) {

	var encrypted, resBt []byte
	if encrypted, err = base64.StdEncoding.DecodeString(text); err != nil {
		return
	}
	// 传入密文和自己定义的密钥，需要和加密的密钥一样，不一样会报错 可以自己传入初始化向量,如果不传就使用默认的初始化向量,16字节
	if resBt, err = aes.AesCtrDecrypt(encrypted, []byte(aesKey), []byte(ChatEncryptionKey)); err != nil {
		return
	}
	if resBt, err = pkcs7UnPadding(resBt); err != nil {
		return
	}
	res = string(resBt)
	res = strings.ReplaceAll(res, "\b", "")
	return
}
