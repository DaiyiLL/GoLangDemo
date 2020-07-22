package main

import (
	"fmt"
	"net/http"
)

func netHandler(resp http.ResponseWriter, req *http.Request)  {
	resp.Write([]byte("this is a Web server"))

	fmt.Println("Header:", req.Header)
	fmt.Println("URL:", req.URL)
	fmt.Println("Method:", req.Method)
	fmt.Println("Host:", req.Host)
	fmt.Println("RemoteAddr:", req.RemoteAddr)
	fmt.Println("Body:", req.Body)
}

func main()  {
	// 注册回调函数，该函数在客户端访问服务器是，会自动被调用
	http.HandleFunc("/itcast", netHandler)
	// 绑定服务器监听地址
	http.ListenAndServe("127.0.0.1:8000", nil)
}