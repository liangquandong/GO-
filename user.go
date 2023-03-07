package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

// 用户上线封装
func (s *User) Online() {
	s.server.mapLock.Lock()
	s.server.OnlineMap[s.Name] = s
	s.server.mapLock.Unlock()
	s.server.BroadCast(s, "已上线")

}

// 业务下线封装
func (s *User) Offline() {
	s.server.mapLock.Lock()
	delete(s.server.OnlineMap, s.Name)
	s.server.mapLock.Unlock()
	s.server.BroadCast(s, "下线")
}

// 给客户端发送消息
func (s *User) SendMsg(msg string) {
	s.conn.Write([]byte(msg))
}

// 用户业务封装
func (s *User) UserMessage(msg string) {
	if msg == "who" {
		s.server.mapLock.Lock()
		for _, v := range s.server.OnlineMap {
			onlineMsg := "用户【" + v.Name + "】在线...\n"
			s.SendMsg(onlineMsg)
		}
		s.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		_, ok := s.server.OnlineMap[newName]
		if ok {
			s.SendMsg("该用户名已被使用")
		} else {
			s.server.mapLock.Lock()
			delete(s.server.OnlineMap, s.Name)
			s.server.OnlineMap[newName] = s
			s.server.mapLock.Unlock()
			s.Name = newName
			s.SendMsg("你已修改用户名：" + newName + "\n")
		}

	} else if len(msg) > 4 && msg[:3] == "to|" {
		privateName := strings.Split(msg, "|")[1]
		if privateName == "" {
			s.SendMsg("格式错误，私聊请使用“to|私聊对象名称|消息”格式")
			return
		}
		//查看私聊对象是否存在或下线
		privateUser, ok := s.server.OnlineMap[privateName]
		if !ok {
			s.SendMsg("用户不存在或已下线")
			return
		}
		//禁止发空消息
		privateMsg := strings.Split(msg, "|")[2]
		if privateMsg == "" {
			s.SendMsg("消息不能为空")
			return
		}
		//发送消息给对方
		privateUser.SendMsg("【" + s.Name + "】" + "用户：" + privateMsg)

	} else {
		s.server.BroadCast(s, msg)
	}
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n")) //给正在连接的客户端发送消息
	}
}
