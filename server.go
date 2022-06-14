package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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

	user := NewUser(conn, this)

	// 监听用户是否活跃
	isLive := make(chan bool)
	
	user.Online()

	// receive infomation from user
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

			// get user's info without "\n"
			msg := string(buf[:n-1])
			// resovle msg
			user.DoMessage(msg)
			// a user is line if he send a message 
			isLive <- true
		}
	}()

	// block handler to avoid delete user
	for{
		select {
			case <- isLive:
				//当前用户是活跃的，应该重置定时器
				//不做任何事情，为了激活select，更新下面的定时器

			case <-time.After(time.Second * 300):
				// 已经超时，强制关闭当前User
				user.SendMsg("您超时了，已被移出当前聊天室！\n")

				// 销毁资源
				close(user.C)

				// 关闭连接
				conn.Close()

				// 退出当前Handler
				return
		}
	}
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
