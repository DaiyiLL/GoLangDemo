package main

import (
	"fmt"
	"strings"
)

// 测试文件
func main()  {
	//TestDes()
	//TestAes()

	//TestRsa()

	MyHash()
}

func TestRsa()  {
	//GenerateRsaKey(1024)

	plainText := []byte("daishuyi123456")
	cipherText := RSAEncrypt(plainText, "public.pem")
	fmt.Println(string(cipherText))
	fmt.Println("=========================================")
	plainText = RSADecrypt(cipherText, "privateKey.pem")
	fmt.Println(string(plainText))
}


func TestAes()  {
	plainText := []byte("dsy1672dsy3821648732")
	key := []byte(strings.Repeat("1234", 4))
	cipherText := AesEncrypt(plainText, key)
	fmt.Println(string(cipherText))
	plainText1 := AesEncrypt(cipherText, key)
	fmt.Println(string(plainText1))
}

func TestDes()  {
	plainText := []byte("dsy1672dsy3821648732")
	key := []byte("1234abcd")
	cipherText := DesEncrypt(plainText, key)
	fmt.Println(string(cipherText))
	plainText1 := DesDecrypt(cipherText, key)
	fmt.Println(string(plainText1))
}
