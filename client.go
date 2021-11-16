package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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

//deal with server message, directly stdout
func (client *Client) DealResponse() {
	//once client.conn has message, copy to standard output, block listen
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1 Public chat")
	fmt.Println("2 Pravite chat")
	fmt.Println("3 Update name")
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

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println("* Please enter a username, type \"exit\" for exit *")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("* Please enter the content, type \"exit\" for exit *")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			//消息不为空则发送
			if len(chatMsg) != 0 {
				sendMsg := "to " + remoteName + " " + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println("* Please enter the content, type \"exit\" for exit *")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUsers()
		fmt.Println("* Please enter a username, type \"exit\" for exit *")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) PublicChat() {
	//Hint user to tpye in message
	var chatMsg string

	fmt.Println("* Please enter the content, type \"exit\" for exit *")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		//send to if message is not empty
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println("* Please enter the content, type \"exit\" for exit *")
		fmt.Scanln(&chatMsg)
	}

}

func (client *Client) UpdateName() bool {

	fmt.Println("* Type the new name: *")
	fmt.Scanln(&client.Name)

	sendMsg := "rename " + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		//according to user's choice
		switch client.flag {
		case 1:
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
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

	//create a goroutine to deal with server message
	go client.DealResponse()

	fmt.Println("* Success to connect to server *")

	//start client side
	client.Run()

}
