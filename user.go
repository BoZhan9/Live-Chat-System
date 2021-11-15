package main

import "net"

type User struct {
	Name string
	Addr string
	C chan string //check if there is a message
	conn net.Conn //socket connection

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	//get address
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
		server: server,
	}
	//start a tap goroutine
	go user.ListenMessage()

	return user
}

//user online
func (t *User) Online() {

	//user online, add into map
	t.server.mapLock.Lock()
	t.server.OnlineMap[t.Name] = t
	t.server.mapLock.Unlock()

	//broadcast online
	t.server.BroadCast(t, "* Enter the chat *")
}

//user offline
func (t *User) Offline() {

	//user offline, delete from map
	t.server.mapLock.Lock() //map thread-unsafe, add lock 
	delete(t.server.OnlineMap, t.Name)
	t.server.mapLock.Unlock()

	//broadcast offline
	t.server.BroadCast(t, "* Leave the chat *")

}

//user send message
func (t *User) DoMessage(msg string) {
	t.server.BroadCast(t, msg)
}

//tap current user channel, once got message, directly send to user
func (t *User) ListenMessage() {
	for {
		msg := <-t.C //read connection message

		t.conn.Write([]byte(msg + "\n")) //convert to binary array
	}
}