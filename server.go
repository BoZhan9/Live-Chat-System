package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip string
	Port int

	//online user map
	OnlineMap map[string]*User
	mapLock sync.RWMutex //add a lock
	
	//message broadcast channel
	Message chan string
}

//construct server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}

//tap message broadcast channel's goroutine
//once got message send to all online users
func (t *Server) ListenMessager() {
	for {
		msg := <-t.Message

		t.mapLock.Lock()
		//for loop value, which are user obj 
		for _, cli := range t.OnlineMap { 
			cli.C <- msg
		}
		t.mapLock.Unlock()
	}
}

func (t *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	t.Message <- sendMsg
}

func (t *Server) Handler(conn net.Conn) {
	//...tasks for current connection
	//fmt.Println("Connection successful")
	user := NewUser(conn, t)

	user.Online()

	//get user message
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err", err)
				return
			}

			//get user message delete '\n'
			msg := string(buf[:n-1])

			//user send message
			user.DoMessage(msg)
		}
	}()

	//set process state as blocked
	select {}
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

	//start a tap goroutine
	go t.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		//open a goroutine to do handler
		go t.Handler(conn)
	}
}