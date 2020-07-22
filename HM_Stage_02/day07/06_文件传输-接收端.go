package main

import (
	"fmt"
	"net"
	"os"
)

func main()  {
	// 创建用于监听的socket
	listener, err := net.Listen("tcp", "127.0.0.1:8008")
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	defer listener.Close()

	// 阻塞监听
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("listener.Accept err: ", err)
		return
	}
	defer conn.Close()

	// 获取文件名，保存
	buf := make([]byte, 4096)
	// 文件名的长度不可以超过1024个字节
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read err: ", err)
		return
	}
	fileName := string(buf[:n])
	fmt.Println("文件名: ", fileName)
	// 回写ok给发送端
	conn.Write([]byte("ok"))

	// 获取文件内容
	recvFile(conn, fileName)
}

func recvFile(conn net.Conn, fileName string)  {
	// 创建文件 按照文件名创建新文件
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println("os.Create err:", err)
		return
	}
	defer f.Close()
	// 从网络中毒数据，写入本地文件中
	buf := make([]byte, 4096)
	for {
		n, _ := conn.Read(buf)
		if n == 0 {
			fmt.Println("接收文件完毕")
			return
		}
		// 写入本地文件，
		f.Write(buf[:n])
	}
}
