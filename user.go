package main

import "net"

type User struct {
	Name	string
	Addr	string
	C 		chan string
	conn	net.Conn
}

// 创建用户API
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
	}
	// 启动监听协程
	go user.ListenMessage()
	return user
}

// 监听当前User channel的方法，一旦有消息，直接发给对应的客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}