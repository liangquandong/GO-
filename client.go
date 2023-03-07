package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

// 连接服务器对象
func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	//连接服务端
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverPort, serverPort))
	if err != nil {
		fmt.Println("net.Dial:", err)
		return nil
	}
	client.conn = conn
	return client
}

func main() {
	client := NewClient("127.0.0.1", 9999)
	if client == nil {
		fmt.Println(">>>>>服务器连接失败<<<<<")
		return
	}
	select {}
}
