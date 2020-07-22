package main

import (
	"fmt"
	"regexp"
)

func main()  {
	str := "abc a7c mfc cat 8ca azc cba"
	// 解析，编译正则表达式
	ret := regexp.MustCompile(`a.c`)  // ``: 表示使用原生字符串
	// 提取，需要的信息
	alls := ret.FindAllString(str, -1)
	fmt.Println(alls)
}
