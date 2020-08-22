package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
)

const (
	CIPHER_IV_TEXT = "12345678"
	CIPHER_IV_SEED = "12345678abcdefgh"
)

// des加密

func PaddingLastGroup(plainText []byte, blockSize int) []byte {
	// 1.求出最后一个组中剩余的字节数
	padNum := blockSize - len(plainText) % blockSize
	// 2. 创建新的切片，长度 == padNum, 每个字节值为 byte(padNum)
	char := []byte{(byte(padNum))}
	// 切片创建，并初始化
	newPlaint := bytes.Repeat(char, padNum)
	// 3. newPlain数组追加到原始明文后面
	plainText = append(plainText, newPlaint...)
	return plainText
}

// 去掉填充的数据
func unPaddingLastGroup(plainText []byte) []byte  {
	// 1. 获取最后一个字节
	length := len(plainText)
	lastChar := plainText[length - 1]
	number := int(lastChar)  // 尾部填充的字节个数
	return plainText[:length - number]
}

// des加密
func DesEncrypt(plainText, key []byte) []byte  {
	// 1. 创建一个底层使用des的密码接口
	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// 2. 明文填充
	newText := PaddingLastGroup(plainText, block.BlockSize())

	// 3. 创建一个使用cbc分组接口
	iv := []byte(CIPHER_IV_TEXT)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	// 4. 加密
	cipherText := make([]byte, len(newText))
	blockMode.CryptBlocks(cipherText, newText)
	return cipherText
}

// des解密
func DesDecrypt(cipherText, key []byte) []byte {
	// 1. 创建一个底层使用des的密码接口
	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// 2. 创建一个使用cbc分组模式解密的接口
	iv := []byte(CIPHER_IV_TEXT)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	// 3. 解密(将所有的密文还原为明文)
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	// 4. 明文去除填充
	return unPaddingLastGroup(plainText)
}

// Aes加密, 分组模式ctr计数器
func AesEncrypt(plainText, key []byte) []byte  {
	// 1. 创建一个底层使用aes的密码接口
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//// 2. 明文填充
	//newText := PaddingLastGroup(plainText, block.BlockSize())

	// 3. 创建一个使用ctr分组接口
	iv := []byte(CIPHER_IV_SEED)
	stream := cipher.NewCTR(block, iv)
	// 4. 加密
	cipherText := make([]byte, len(plainText))
	stream.XORKeyStream(cipherText, plainText)
	return cipherText
}

// des解密
func AesDecrypt(cipherText, key []byte) []byte {
	// 1. 创建一个底层使用aes的密码接口
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// 2. 创建一个使用cbc分组模式解密的接口
	iv := []byte(CIPHER_IV_SEED)
	stream := cipher.NewCTR(block, iv)
	// 3. 解密(将所有的密文还原为明文)
	plainText := make([]byte, len(cipherText))
	stream.XORKeyStream(plainText, cipherText)
	return plainText
}