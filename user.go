package main

import "net"

type User struct {
	Name	string
	Addr	string
	C 		chan string
	conn	net.Conn
	server  *Server
}

// 创建用户API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
		server: server,
	}
	// 启动监听协程
	go user.ListenMessage()
	return user
}

// 用户上线业务
func (this *User) Online() {
	// 用户上线，将用户添加到OnlineMap
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 广播当前用户上线消息
	this.server.Broadcast(this, "上线啦!")
}

// 用户下线业务
func (this *User) Offline() {
	// 用户下线，将用户信息从OnlineMap中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 广播当前用户上线消息
	this.server.Broadcast(this, "回家吃饭啦!")
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	this.server.Broadcast(this, msg)
}


// 监听当前User channel的方法，一旦有消息，直接发给对应的客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}