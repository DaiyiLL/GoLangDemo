package main

import (
	"fmt"
	"net"
	"time"
)

func main()  {
	// 组织一个 udp 地址结构
	srvAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8003")
	if err != nil {
		fmt.Println("net.ResolveIPAddr err: ", err)
		return
	}

	fmt.Println("UDP 服务器地址结构，创建完成!!!!")

	// 创建用户通信的socket
	udpConn, err := net.ListenUDP("udp", srvAddr)
	if err != nil {
		fmt.Println("net.ListenUDP err: ", err)
		return
	}
	defer udpConn.Close()
	fmt.Println("UDP 服务器通信socket，创建完成!!!!")

	// 读取客户端发送的数据
	buf := make([]byte, 4096)
	// 返回3个值，分别是读取到的字节数，客户端的地址，err
	n, cltAddr, err := udpConn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("udpConn.ReadFromUDP err: ", err)
		return
	}
	// 模拟处理数据
	fmt.Printf("服务器读到 %v 的数据 : %s\n", cltAddr, string(buf[:n]))
	// 回写数据给客户端
	daytime := time.Now().String()
	_, err = udpConn.WriteToUDP([]byte(daytime), cltAddr)
	if err != nil {
		fmt.Println("udpConn.WriteToUDP err: ", err)
		return
	}
}
