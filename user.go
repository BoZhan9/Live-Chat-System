package main

import (
	"net"
	"strings"
)

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

func (t *User) SendMsg(msg string) {
	t.conn.Write([]byte(msg))
}

//user send message
func (t *User) DoMessage(msg string) {
	//search online user
	if msg == "who" {
		t.server.mapLock.Lock()
		for _, user := range t.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "* Online *\n"
			t.SendMsg(onlineMsg)
		}
		t.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename " {
		//format: rename Brian
		newName := strings.Split(msg, " ")[1]

		//check if name already exist
		_, ok := t.server.OnlineMap[newName]
		if ok {
			t.SendMsg("* Name already exist *\n")
		} else {
			t.server.mapLock.Lock()
			delete(t.server.OnlineMap, t.Name)
			t.server.OnlineMap[newName] = t
			t.server.mapLock.Unlock()

			t.Name = newName
			t.SendMsg("* You have updated name: " + t.Name + " *\n")
		}

	} else if len(msg) > 4 && msg[:3] == "to " {
		//format:  to Brian Hello

		//get receiver username
		remoteName := strings.Split(msg, " ")[1]
		if remoteName == "" {
			t.SendMsg("* Format error, please use \"to Brian Hello\" *\n")
			return
		}

		//get user object by username
		remoteUser, ok := t.server.OnlineMap[remoteName]
		if !ok {
			t.SendMsg("* The user you want to send doesn't exist *\n")
			return
		}

		//get message content, send to user obkect
		content := strings.Split(msg, " ")[2]
		if content == "" {
			t.SendMsg("* No content detect, type again *\n")
			return
		}
		remoteUser.SendMsg(t.Name + " said to you: " + content)

	} else {
		t.server.BroadCast(t, msg)
	}
}

//tap current user channel, once got message, directly send to user
func (t *User) ListenMessage() {
	for {
		msg := <-t.C //read connection message

		t.conn.Write([]byte(msg + "\n")) //convert to binary array
	}
}