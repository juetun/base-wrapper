package common

import (
	"encoding/base64"
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
	var cryptText []byte

	if cryptText, err = aes.AesCtrEncrypt([]byte(text), []byte(aesKey), []byte(ChatEncryptionKey)); err != nil {
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
	if resBt, err = aes.AesCtrDecrypt(encrypted, []byte(aesKey), []byte(ChatEncryptionKey)); err != nil {
		return
	}
	res = string(resBt)

	res = strings.ReplaceAll(res, "", "")
	res = strings.ReplaceAll(res, "\b", "")
	return
}
