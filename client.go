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
	flag int
}

func NewClient(serverIp string, serverPort int) *Client {
	//construct client object
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
		flag: 999,
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

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1 Public chat")
	fmt.Println("2 Pravite chat")
	fmt.Println("3 Rename")
	fmt.Println("0 Exit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("* Please enter a valid number (0-3) *")
		return false
	}

}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		//according to user's choice
		switch client.flag {
		case 1:
			fmt.Println("Public chat")
			break
		case 2:
			fmt.Println("Private chat")
			break
		case 3:
			fmt.Println("rename")
			break
		}
	}
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
	client.Run()

}
