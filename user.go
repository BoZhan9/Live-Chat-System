package main

import "net"

type User struct {
	Name string
	Addr string
	C chan string //check if there is a message
	conn net.Conn //socket connection
}

func NewUser(conn net.Conn) *User {
	//get address
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
	}
	//start a tap goroutine
	go user.ListenMessage()

	return user
}

//tap current user channel, once got message, directly send to client side
func (t *User) ListenMessage() {
	for {
		msg := <-t.C //read connection message

		t.conn.Write([]byte(msg + "\n")) //convert to binary array
	}
}