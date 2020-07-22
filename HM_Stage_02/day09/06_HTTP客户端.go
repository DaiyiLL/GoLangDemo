package main

import (
	"fmt"
	"net/http"
)

func main() {
	//resp, err := http.Get("https://www.itcast.cn/")
	resp, err := http.Get("http://www.baidu.com/")
	if err != nil {
		fmt.Println("http.Get err: ", err)
		return
	}
	defer resp.Body.Close()

	// 简单查看应答包
	fmt.Println("Header:", resp.Header)
	fmt.Println("StatusCode:", resp.StatusCode)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Proto:", resp.Proto)

	buf := make([]byte, 4096)
	var result string

	for {
		n, err := resp.Body.Read(buf)
		if n == 0 {
			fmt.Println("-------------- Read finish-----------------")
			break
		}
		if err != nil {
			fmt.Println("resp.Body.Read err: ", err)
			break
		}
		result += string(buf[:n])
	}
	fmt.Printf("|%v|\n", result)
}