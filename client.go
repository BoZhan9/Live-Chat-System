package main

import (
	"flag"
	"fmt"
	//"io"
	"net"
	//"os"
)

type Client struct {
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
	//flag int //current client mode
}

func NewClient(serverIp string, serverPort int) *Client {
	//construct client object
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
		//flag: 999,
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

// //deal response from server display on standard output
// func (client *Client) DealResponse() {
// 	//once client.conn has message, copy to stdout, block listen
// 	io.Copy(os.Stdout, client.conn)
// }

// func (client *Client) menu() bool {
// 	var flag int

// 	fmt.Println("1. Public Chat")
// 	fmt.Println("2. Private Chat")
// 	fmt.Println("3. Rename")
// 	fmt.Println("0. Exit")

// 	fmt.Scanln(&flag)

// 	if flag >= 0 && flag <= 3 {
// 		client.flag = flag
// 		return true
// 	} else {
// 		fmt.Println("* Please enter valid number (0-3) *")
// 		return false
// 	}

// }

// //search online client
// func (client *Client) SelectUsers() {
// 	sendMsg := "who\n"
// 	_, err := client.conn.Write([]byte(sendMsg))
// 	if err != nil {
// 		fmt.Println("conn Write err:", err)
// 		return
// 	}
// }

// //private chat
// func (client *Client) PrivateChat() {
// 	var remoteName string
// 	var chatMsg string

// 	client.SelectUsers()
// 	fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
// 	fmt.Scanln(&remoteName)

// 	for remoteName != "exit" {
// 		fmt.Println(">>>>请输入消息内容, exit退出:")
// 		fmt.Scanln(&chatMsg)

// 		for chatMsg != "exit" {
// 			//消息不为空则发送
// 			if len(chatMsg) != 0 {
// 				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
// 				_, err := client.conn.Write([]byte(sendMsg))
// 				if err != nil {
// 					fmt.Println("conn Write err:", err)
// 					break
// 				}
// 			}

// 			chatMsg = ""
// 			fmt.Println(">>>>请输入消息内容, exit退出:")
// 			fmt.Scanln(&chatMsg)
// 		}

// 		client.SelectUsers()
// 		fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
// 		fmt.Scanln(&remoteName)
// 	}
// }

// func (client *Client) PublicChat() {
// 	//提示用户输入消息
// 	var chatMsg string

// 	fmt.Println(">>>>请输入聊天内容，exit退出.")
// 	fmt.Scanln(&chatMsg)

// 	for chatMsg != "exit" {
// 		//发给服务器

// 		//消息不为空则发送
// 		if len(chatMsg) != 0 {
// 			sendMsg := chatMsg + "\n"
// 			_, err := client.conn.Write([]byte(sendMsg))
// 			if err != nil {
// 				fmt.Println("conn Write err:", err)
// 				break
// 			}
// 		}

// 		chatMsg = ""
// 		fmt.Println(">>>>请输入聊天内容，exit退出.")
// 		fmt.Scanln(&chatMsg)
// 	}

// }

// func (client *Client) UpdateName() bool {

// 	fmt.Println("Please enter username:")
// 	fmt.Scanln(&client.Name)

// 	sendMsg := "rename " + client.Name + "\n"
// 	_, err := client.conn.Write([]byte(sendMsg))
// 	if err != nil {
// 		fmt.Println("conn.Write err:", err)
// 		return false
// 	}

// 	return true
// }

// func (client *Client) Run() {
// 	for client.flag != 0 {
// 		for client.menu() != true {
// 		}

// 		//According user entered
// 		switch client.flag {
// 		case 1:
// 			//public chat
// 			client.PublicChat()
// 			break
// 		case 2:
// 			//private chat
// 			client.PrivateChat()
// 			break
// 		case 3:
// 			//rename
// 			client.UpdateName()
// 			break
// 		}
// 	}
// }

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

	//open a single goroutine to deal server side message
	//go client.DealResponse()

	fmt.Println("* Success to connect to server *")

	//start client side
	//client.Run()
	select {}
}
