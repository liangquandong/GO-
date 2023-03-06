package main

import "net"

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

// 用户业务封装
func (s *User) UserMessage(msg string) {
	s.server.BroadCast(s, msg)
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n")) //给正在连接的客户端发送消息
	}
}
