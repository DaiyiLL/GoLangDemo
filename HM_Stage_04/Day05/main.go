package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
)

func main()  {
	//test01()
	test02()
}

func test02()  {
	src := []byte("dajifd我拿到射流风机的数量的数据福利多少")
	//key := []byte("helloworld")
	sigText := SignatureRSA(src, "privateKey.pem")
	flag := VerrityRSA(src, sigText, "public.pem")
	fmt.Println(flag)
}

func test01()  {
	src := []byte("dajifd我拿到射流风机的数量的数据福利多少")
	key := []byte("helloworld")
	hmacText := GenerateHamc(src, key)
	flag := VerifyHamc(src, key, hmacText)
	if flag {
		fmt.Println("消息认证成功")
	} else {
		fmt.Println("消息认证失败")
	}
}

// 生成消息认证码
func GenerateHamc(plainText, key []byte) []byte  {
	// 1. 创建哈希接口，需要制定使用的哈希算法 和 密钥
	myHash := hmac.New(sha1.New, key)
	// 2. 给哈希对象添加数据
	myHash.Write(plainText)
	// 3. 计算散列值
	hashText := myHash.Sum(nil)
	return hashText
}

// 验证消息认证码
func VerifyHamc(plainText, key, hashText[]byte) bool  {
	// 1. 创建哈希接口，需要制定使用的哈希算法 和 密钥
	myHash := hmac.New(sha1.New, key)
	// 2. 给哈希对象添加数据
	myHash.Write(plainText)
	// 3. 计算散列值
	hmacText := myHash.Sum(nil)
	// 4. 两个散列值进行比较
	return hmac.Equal(hashText, hmacText)

	//return false
}
