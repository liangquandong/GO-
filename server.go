package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

/*流程：
1、首先通过NewServer函数定义Server也就是指定ip端口
2、调用Start方法启动服务，同时使用ListenMessager方法监听，并且在Accept阻塞，待有用户上线取消阻塞后调用处理函数Handler
3、Handler函数通过传参net.Conn类型创建User类型的用户表user，并初始化一个OnlineMap，并且使用新建的user用户信息已经msg字符串调用BroadCast方法
4、一旦有用户上线就会触发BroadCast方法，使sendMsg信息添加到Messages的管道里边去，Messages存在数据后ListenMessager就会把数据取出添加到每个用户的C这个管道
5、一旦User下的C有数据可取，就会ListenMessage使其将信息发送给客户端
*/

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Messages  chan string
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Messages:  make(chan string),
	}
	return server
}

// 监听message广播消息的channel的goroutine，一旦有消息就发送给全部的在线user
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Messages
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// 广播消息方法
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Messages <- sendMsg
}

// 处理用户上线函数
func (s *Server) Handler(conn net.Conn) {
	//fmt.Println("连接建立成功......")
	//当连接成功也就是用户上线的时候，调用user.go中的NewUser方法
	user := NewUser(conn, s)
	user.Online()
	//接收客户端的信息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			msg := string(buf[:n-1])
			user.UserMessage(msg)
		}
	}()
	select {}
}

// 启动服务
func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err", err)
		return
	}
	defer listener.Close()
	go s.ListenMessager()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listerner accept err：", err)
			continue
		}
		go s.Handler(conn)
	}
}
