package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	mod        int // 当前client的模式
}

func NewClient(serverIp string, serverPort int) *Client {
	// create the object of client
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		mod:        999,
	}

	// link server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn

	// return object
	return client
}

// resovle the response of server and display that on stdout
func (client *Client) DealResponse() {
	// if there is info in client.conn, then display that on stdout, block to listen forever
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var mod int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&mod)
	if mod >= 0 && mod <= 3 {
		client.mod = mod
		return true
	} else {
		fmt.Println(">>>>请输入合法范围内的数字<<<<")
		return false
	}
}

// querry online user
func (client *Client) QuerryUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write erorr:", err)
		return
	}
}

// private chat
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.QuerryUsers()
	fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
	fmt.Scanln(&remoteName)
	for remoteName != "exit" {
		fmt.Println(">>>>请输入消息内容, exit退出:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write erorr：", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>>请输入消息内容, exit退出:")
			fmt.Scanln(&chatMsg)
		}

		client.QuerryUsers()
		fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>>请输入聊天内容，exit退出:")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write erorr:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>请输入聊天内容，exit退出:")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>请输入用户名:")
	var NewName string
	fmt.Scanln(&NewName)
	sendMsg := "rename|" + NewName + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write erorr:", err)
		return false
	}
	client.Name = NewName
	return true
}

func (client *Client) Run() {
	for client.mod != 0 {
		for client.menu() != true {
		}

		// resovle different job according to different mod
		switch client.mod {
		case 1:
			client.PublicChat()
		case 2:
			client.PrivateChat()
		case 3:
			client.UpdateName()
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

func main() {
	// commanf line parse
	flag.Parse()

	clien := NewClient(serverIp, serverPort)
	if clien == nil {
		fmt.Println(">>>>> 链接服务器失败...")
		return
	}

	// create a goroutine to resovle the response of server
	go clien.DealResponse()

	fmt.Println(">>>>>链接服务器成功...")
	// run client

	clien.Run()
}
