package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)



// 1. 生成ras的秘钥对，并且保存到磁盘文件中
func GenerateRsaKey(keySize int) error {
	// 1. 使用rsa.GenerateKey生成密钥
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		fmt.Println("rsa.GenerateKey: ", err)
		return err
	}
	// 2. 通过x509标准将得到的rsa私钥序列化为ASN.1 的DER编码字符串
	derText := x509.MarshalPKCS1PrivateKey(privateKey)
	// 3. 要组织一个pem.block
	block := pem.Block{
		Type: "rsa private key",
		Bytes: derText,
	}
	// 4. pem编码
	file, err := os.Create("privateKey.pem")
	if err != nil {
		fmt.Println("os.Create: ", err)
		return err
	}
	err = pem.Encode(file, &block)
	if err != nil {
		fmt.Println("pem.Encode: ", err)
		return err
	}
	file.Close()

	// ================================ 公钥 ==========================
	// 1. 从私钥中取出公钥
	publicKey := privateKey.PublicKey
	// 2. 使用x509标准进行格式化  （取地址）
	derstream, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		fmt.Println("x509.MarshalPKIXPublicKey: ", err)
		return err
	}
	// 3. 将得到的数据放到pem.Block中
	publicBlock := pem.Block{
		Type: "rsa public key",
		Bytes: derstream,
	}
	// 4. pem编码
	pubFile, err := os.Create("public.pem")
	pem.Encode(pubFile, &publicBlock)
	if err != nil {
		fmt.Println("public pem.Encode: ", err)
		return err
	}
	pubFile.Close()

	return nil
}

// 1. 生成ras的秘钥对，并且保存到磁盘文件中
func GenerateRsaPublicKey(keySize int) error {
	// 1. 使用rsa.GenerateKey生成密钥
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		fmt.Println("rsa.GenerateKey: ", err)
		return err
	}
	// 2. 通过x509标准将得到的rsa私钥序列化为ASN.1 的DER编码字符串
	derText := x509.MarshalPKCS1PrivateKey(privateKey)
	// 3. 要组织一个pem.block
	block := pem.Block{
		Type: "rsa private key",
		Bytes: derText,
	}
	// 4. pem编码
	file, err := os.Create("privateKey.pem")
	if err != nil {
		fmt.Println("os.Create: ", err)
		return err
	}
	err = pem.Encode(file, &block)
	if err != nil {
		fmt.Println("pem.Encode: ", err)
		return err
	}

	return nil
}

// RSA加密
func RSAEncrypt(plaintText []byte, fileName string) []byte {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	buf := make([]byte, fileInfo.Size())
	_, err = file.Read(buf)
	if err != nil {
		panic(err)
	}
	file.Close()
	// 2. pem解码
	block, _ := pem.Decode(buf)
	if err != nil {
		panic(err)
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if ok {
		// 使用公钥加密
		cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, plaintText)
		if err != nil {
			panic(err)
		}
		return cipherText
	}
	return nil
}

// RSA解密
func RSADecrypt(cipherText []byte, fileName string) []byte {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	buf := make([]byte, fileInfo.Size())
	_, err = file.Read(buf)
	if err != nil {
		panic(err)
	}
	file.Close()
	// 2. pem解码
	block, _ := pem.Decode(buf)
	if err != nil {
		panic(err)
	}


	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	// 使用私钥加密
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		panic(err)
	}
	return plainText
}
