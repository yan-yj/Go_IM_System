package main

import (
	"net"
	"strings"
)

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

	// 发送提示性信息
	info := "您已成功登录！\n"
	this.SendMsg(info)
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

// 给当前User发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线用户
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name +":" + "在线等你哦~\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
		
	} else if len(msg) > 7 && msg[:7] == "rename|" { 
		// 消息格式rename|XX
		newName := strings.Split(msg, "|")[1]
		
		// 判断用户名是否已经存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("当前用户名已被使用")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMsg("您已经成功更新用户名为:" + this.Name +"\n")
		}
		
	} else if len(msg) > 4 && msg[:3] == "to|"{
		// 消息格式："to|张三|内容"
		// 1 从msg中获取私聊对象用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("消息格式有误，请使用“to|张三|内容”发送消息\n")
			return
		}
		// 2 根据用户名获取对方User对象
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("你输入的用户不存在，请确认无误后再发\n")
			return
		}

		// 3 通过User对象将消息发送过去
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("消息内容为空，请确认无误后重发")
			return
		}
		remoteUser.SendMsg(this.Name + "对您说：" + content + "\n")

	} else {
		this.server.Broadcast(this, msg)
	}
}


// 监听当前User channel的方法，一旦有消息，直接发给对应的客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}