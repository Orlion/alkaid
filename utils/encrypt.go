package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
)

//输入明文，输出密文
func AesEncrypt(plainText string, key string) (string, error) {
	// 第一步：创建aes密码接口
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	//第二步：创建分组模式ctr
	// iv 要与算法长度一致，16字节
	// 使用bytes.Repeat创建一个切片，长度为blockSize()，16个字符"1"
	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)

	//第三步：加密
	dst := make([]byte, len(plainText))
	stream.XORKeyStream(dst, []byte(plainText))

	return string(dst), nil
}

//输入密文，得到明文
func AesDecrypt(encryptData string, key string) (string, error) {
	return AesEncrypt(encryptData, key)
}

func Md5(plainText string) string {
	h := md5.New()
	h.Write([]byte(plainText))
	return hex.EncodeToString(h.Sum(nil))
}
