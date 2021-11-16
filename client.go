package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	//construct client object
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
	}

	//connect to server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn

	return client
}

var serverIp string
var serverPort int

//./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "* Set IP address *")
	flag.IntVar(&serverPort, "port", 8888, "* Set port *")
}

func main() {
	//parse commend
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("* Fail to connect to server *")
		return
	}

	fmt.Println("* Success to connect to server *")

	//start client side
	//client.Run()
	select {}
}
