package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip 	 string
	Port int

	// the list of online user
	OnlineMap 	map[string]*User
	mapLock 	sync.RWMutex   // 读写锁

	// channel that can be used to broadcast
	Message		chan string
}

// create a interface of server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip: 		ip,
		Port: 		port,
		OnlineMap: 	make(map[string]*User),
		Message: 	make(chan string),
	}

	return server
}

// send message to the user who is online
func (this *Server) ListenMessager() {
	for {
		msg := <- this.Message

		// visit and send msg to all online user
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// broadcast message
func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

// do the work of this conn
func (this *Server) handler(conn net.Conn) {
	fmt.Println("Create the conn successflly!")

	user := NewUser(conn)
	// 用户上线，将用户添加到OnlineMap
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 广播当前用户上线消息
	this.Broadcast(user, "上线啦！")

	// block handler to avoid delete user
	select {}
}

// the interface of starting server
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer listener.Close()

	// start listenMesager
	go this.ListenMessager()

	for {
		// accepet
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		// do handler
		go this.handler(conn)
	}
}
