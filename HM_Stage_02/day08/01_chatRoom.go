package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// 消息实体
type Message struct {
	From string
	Msg  string
}

// 创建用户题结构体
type Client struct {
	C    chan Message
	Name string
	Addr string
}

// 创建全局map，存储在线用户
var onlineMap map[string]Client
// 创建全局channel传递用户消息
var message = make(chan Message)



func main()  {

	// 创建监听套接字
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	defer listener.Close()

	// 创建管理者go程，管理map 和全局channel
	go Manager()

	// 循环监听客户端连接请求
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err: ", err)
			return
		}
		go HandlerConnect(conn)
	}
}

func Manager()  {
	// 初始化在线用户 onlineMap
	onlineMap = make(map[string]Client)

	// 监听全局channel中是否有数据, 没有数据的时候阻塞
	for {
		msg := <- message

		// 循环发送消息给 所有在线用户
		for _, client := range onlineMap {
			client.C <- msg
		}
	}

}

func HandlerConnect(conn net.Conn)  {
	defer conn.Close()

	// 创建channel判断用户是否活跃
	hasData := make(chan bool)

	// 创建新连接用户的结构体
	// 获取用户的 网络地址
	netAddr := conn.RemoteAddr().String()

	client := Client{make(chan Message), netAddr, netAddr}

	// 将新建连接用户，添加到在线用户map中
	onlineMap[netAddr] = client

	// 创建专门用来给当前 用户发送消息的go程
	go WriteMsgToClient(client, conn)

	message <- MakeMsg(client, "login")

	// 创建一个channel，用来判断用户退出状态
	isQuit := make(chan bool)

	// 创建一个匿名go程，专门处理用户发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				fmt.Printf("检测到客户端: %s退出\n", client.Name)
				isQuit <- true
				return
			}
			if err != nil {
				fmt.Println("conn.Read err:", err)
				return
			}
			msg := string(buf[:n - 1])
			if msg == "who" && len(msg) == 3 {
				// 提取在线用户列表
				conn.Write([]byte("online user list:\n"))
				// 遍历当前map，获取在线用户
				for _, user := range  onlineMap {
					userInfo := user.Addr + ":" + user.Name + "\n"
					conn.Write([]byte(userInfo))
				}
			}  else if len(msg) >= 8 && msg[:7] == "rename|" {
				//name := strings.Split(msg, "|")[1]  // msg[8:]
				name := msg[7:]
				client.Name = name  //修改结构体成员name
				onlineMap[netAddr] = client

				conn.Write([]byte("rename successful!\n"))
			} else {
				// 将读到的用户消息，写入到message中
				message <- MakeMsg(client, msg)
			}
			hasData <- true
		}
	}()

	// 保证 不退出
	for {
		// 监听channel上的数据流动
		select {
		case <-isQuit:
			close(client.C)
			delete(onlineMap, client.Addr)
			message <- MakeMsg(client, " logout")   // 写入用户退出消息到全局message
			return
			case <-hasData:
		case <-time.After(time.Second * 10):
			close(client.C)
			delete(onlineMap, client.Addr)
			message <- MakeMsg(client, " time out leaved")   // 写入用户退出消息到全局message
			return
		}
	}
}

func WriteMsgToClient(clnt Client, conn net.Conn)  {
	// 监听 用户自带channel
	for msg := range clnt.C {
		if msg.From != clnt.Name {
			if strings.HasSuffix(msg.Msg,"\n") {
				conn.Write([]byte(msg.Msg))
			} else {
				conn.Write([]byte(msg.Msg + "\n"))
			}
		}
	}
}

func MakeMsg(client Client, msg string) Message  {
	return Message{client.Name, "[" + client.Addr + "]" + client.Name + ": " + msg}
}
