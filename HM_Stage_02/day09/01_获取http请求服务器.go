package main

import (
	"fmt"
	"net"
	"os"
)

func ErrFunc(err error, info string)  {
	if err != nil {
		fmt.Println(info, err)
		//runtime.Goexit()    // 结束当前go程
		os.Exit(-1)      // 将当前进程退出
	}
}

func main0101()  {
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	ErrFunc(err, "net.Listen Err:")

	defer listener.Close()

	conn, err := listener.Accept()
	ErrFunc(err, "Accept err:")
	defer conn.Close()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if n == 0 {
		return
	}
	ErrFunc(err, "Read Err")

	fmt.Println(string(buf[:n]))
}
