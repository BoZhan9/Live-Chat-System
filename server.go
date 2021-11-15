package main

import (
	"fmt"
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

		//将msg发送给全部的在线User
		t.mapLock.Lock()
		for _, cli := range t.OnlineMap { //for loop value, which are client obj 
			cli.C <- msg
		}
		t.mapLock.Unlock()
	}
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (t *Server) Handler(conn net.Conn) {
	//...tasks for current connection
	//fmt.Println("Connection successful")
	user := NewUser(conn)

	//when user is online, add to map
	t.mapLock.Lock() //map thread-unsafe, add lock 
	t.OnlineMap[user.Name] = user
	t.mapLock.Unlock()

	t.BroadCast(user, " is online")

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