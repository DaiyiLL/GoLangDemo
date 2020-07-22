package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main()  {
	// 指定服务器的 ip + port 创建 通信套接字
	conn, err := net.Dial("tcp", "127.0.0.1:8008")
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return
	}
	defer conn.Close()

	// 主动写数据给服务器
	conn.Write([]byte("zhaogong.ipa"))

	// 接收服务器回收的数据
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)

	if err != nil {
		fmt.Println("conn.Read error:", err)
		return
	}

	result := string(buf[:n])
	if result == "ok" {
		sendFile(conn, "/Users/szdjy/Desktop/IPA/ZG/Dis/ExportOptions.plist/zhaogong.ipa")
	}

}

func sendFile(conn net.Conn, filePath string)  {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("os.Open err: ", err)
		return
	}
	defer f.Close()

	buf := make([]byte, 4096)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("文件发送完成")
			} else {
				fmt.Println("f.Read err: ", err)
			}
			return
		}
		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("conn.Write err: ", err)
			return
		}
	}
}