package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip string
	Port int
}

//construct server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
	}
	return server
}

func (t *Server) Handler(conn net.Conn) {
	//...tasks for current connection
	fmt.Println("Connection successful")
}

//start server
func (t *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", t.Ip, t.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close listen socket
	defer listener.Close()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		//do handler
		go t.Handler(conn)
	}
}