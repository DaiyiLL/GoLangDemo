package main

import (
	"fmt"
	"net"
	"os"
)

func ErrFunc1(err error, info string)  {
	if err != nil {
		fmt.Println(info, err)
		//runtime.Goexit()    // 结束当前go程
		os.Exit(-1)      // 将当前进程退出
	}
}

func main()  {
	conn, err := net.Dial("tcp","127.0.0.1:8000")
	ErrFunc1(err, "net.Dial err: ")
	defer conn.Close()

	httpRequest := "GET /itcast HTTP/1.1\r\nHost:127.0.0.1:8000\r\n\r\n"
	conn.Write([]byte(httpRequest))

	buf := make([]byte, 4096)
	n, _ := conn.Read(buf)
	if n == 0 {
		return
	}
	fmt.Printf("|%s|\n", string(buf[:n]))
}
