package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func MyHash()  {
	myHash := sha256.New()
	myHash.Write([]byte("daishuyi12345678"))
	myHash.Write([]byte("Hello world"))
	// 计算结果
	res := myHash.Sum(nil)
	// 格式化
	myStr := hex.EncodeToString(res)
	fmt.Println(myStr)
}
