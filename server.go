package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip 	 string
	Port int
}

// create a interface of server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
	}

	return server
}

// do the work of this conn
func (this *Server) handler(conn net.Conn) {
	fmt.Println("Create the conn successflly!")
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
