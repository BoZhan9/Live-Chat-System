package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

type Client struct {
	ServerIp string
	ServerPort int
	Name string
	conn net.Conn
	flag int
}

func NewClient(serverIp string, serverPort int) *Client {
	// construct client object
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
		flag: 999,
	}

	// connect to server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}

// deal with server message
func (client *Client) DealResponse() {
	// once client.conn has message, copy to standard output, block listen
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var f string

	fmt.Println("")
	fmt.Println("1 Public chat")
	fmt.Println("2 Pravite chat")
	fmt.Println("3 Update name")
	fmt.Println("0 Exit")
	fmt.Println("")

	fmt.Scanln(&f)

	valid := map[string]bool{"0": true, "1": true, "2": true, "3": true}

	if valid[f] {
		flag, _ := strconv.Atoi(f)
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
	fmt.Println("* Username: \"exit\" *")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("* (Private) to " + remoteName + " *")
		fmt.Println("* Don't type in space use \"-\" to replace *")
		fmt.Println("* Key word: \"exit\" *")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			// send if message not empty
			if len(chatMsg) != 0 {
				sendMsg := "to " + remoteName + " " + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println("* (Private) to " + remoteName + " *")
			fmt.Println("* Don't type in space use \"-\" to replace")
			fmt.Println("* Key word: \"exit\" *")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUsers()
		fmt.Println("* Please enter a username *")
		fmt.Println("* Key word: \"exit\" *")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) PublicChat() {
	var chatMsg string

	fmt.Println("* (Public) *")
	fmt.Println("* Don't type in space use \"-\" to replace")
	fmt.Println("* Key word: \"who\" \"exit\" *")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		// send if message is not empty
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println("* (Public) *")
		fmt.Println("* Don't type in space use \"-\" to replace")
		fmt.Println("* Key word: \"who\" \"exit\" *")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) UpdateName() bool {

	fmt.Println("* Type a new name: *")
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
	// parse commend
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("* Fail to connect to server *")
		return
	}
	// create a goroutine to deal with server message
	go client.DealResponse()

	fmt.Println("* Success to connect to server *")
	// start client side
	client.Run()
}
