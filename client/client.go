package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	ImConstant "github.com/LunaY77/IM-go/const"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	command    int
	heartBeat  chan bool
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		command:    999,
		heartBeat:  make(chan bool),
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error", err)
		return nil
	}

	client.conn = conn

	return client
}

func (this *Client) listenMessage() {
	buf := make([]byte, ImConstant.BufLen)
	for {
		n, err := this.conn.Read(buf)
		if n == 0 {
			return
		}
		if err != nil {
			if err == io.EOF {
				fmt.Println("connection closed")
			} else {
				fmt.Println("Conn Read error", err)
			}
			os.Exit(0)
		}

		msg := strings.TrimSpace(string(buf[:n]))
		if msg == ImConstant.IdleMsg {
			fmt.Println("idle... close client connection")
			os.Exit(0)
		}

		fmt.Printf("\n%s\n\n", msg)
	}
}

func (this *Client) Run() {
	for this.command != 0 {
		for this.menu() != true {
		}
		switch this.command {
		case 1:
			// board cast
			this.boardCast()
			break
		case 2:
			// private chat
			this.privateChat()
			break
		case 3:
			// rename
			this.rename()
			break
		case 4:
			// who
			this.who()
			break
		}
	}
}

func (this *Client) menu() bool {
	var command int

	fmt.Println("1. board cast")
	fmt.Println("2. private chat")
	fmt.Println("3. rename")
	fmt.Println("4. who")
	fmt.Println("0. exit")

	_, _ = fmt.Scanln(&command)

	if command >= 0 && command <= 4 {
		this.command = command
		return true
	} else {
		fmt.Println(">>>> please enter a valid number <<<<<<<")
		return false
	}
}

func (this *Client) boardCast() {
	this.sendUtilExit("%s\n")
}

func (this *Client) privateChat() {
	this.who()

	var targetUserName string

	fmt.Println(">>>>> enter target [username]: ")
	_, _ = fmt.Scanln(&targetUserName)

	for targetUserName != "exit" {
		this.sendUtilExit(ImConstant.ToCommand + targetUserName + "|" + "%s\n\n")
		this.who()
		fmt.Println(">>>>> enter target [username]: ")
		_, _ = fmt.Scanln(&targetUserName)
	}
}

func (this *Client) rename() {
	fmt.Println(">>>> enter user name:")
	_, _ = fmt.Scanln(&this.Name)

	sendMsg := ImConstant.RenameCommand + this.Name + "\n"
	this.sendOnce(sendMsg)
}

func (this *Client) who() {
	sendMsg := ImConstant.WhoCommand + "\n"
	this.sendOnce(sendMsg)
}

func (this *Client) sendUtilExit(format string) {
	var content string
	fmt.Println(">>>>> enter content:")
	_, _ = fmt.Scanln(&content)
	for content != "exit" {
		if len(content) != 0 {
			sendMsg := fmt.Sprintf(format, content)
			if !this.sendOnce(sendMsg) {
				return
			}
		}
		content = ""
		fmt.Println(">>>>> enter content:")
		_, _ = fmt.Scanln(&content)
	}
}

func (this *Client) sendOnce(msg string) bool {
	_, err := this.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("conn.Writer error", err)
		return false
	}
	return true
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set server ip")
	flag.IntVar(&serverPort, "port", 8888, "set server port")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>> error connecting server")
		return
	}

	fmt.Println(">>>>> connecting to server")

	go client.listenMessage()
	client.Run()
}
