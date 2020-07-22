package main

import (
	"fmt"
	"net"
)

func main()  {
	// 指定服务器的 ip + port 创建 通信套接字
	conn, err := net.Dial("udp", "127.0.0.1:8003")
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return
	}
	defer conn.Close()

	// 主动写数据给服务器
	conn.Write([]byte("exit"))

	// 接收服务器回收的数据
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)

	if n == 0 {
		fmt.Println("服务器检测到客户端已经管理，断开连接")
		return
	}

	if err != nil {
		fmt.Println("conn.Read error:", err)
		return
	}
	fmt.Println("服务器回发：", string(buf[:n]))
}