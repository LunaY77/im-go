package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	ImConstant "github.com/LunaY77/im-go/const"
)

type Server struct {
	Ip   string
	Port int
	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

func (this *Server) MessageListener() {
	for {
		msg := <-this.Message

		// 将 msg 发送给全部的在线 User
		this.mapLock.Lock()
		for _, client := range this.OnlineMap {
			client.Chan <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) Handler(conn net.Conn) {
	// 心跳
	heartBeat := make(chan bool)
	// 创建 user
	user := NewUser(conn, heartBeat)
	// 用户上线
	user.Online(this)

	// 接收客户端的消息
	go user.Message(this)

	for {
		select {
		case <-heartBeat:
			// 心跳
		case <-time.After(time.Second * ImConstant.IdleTime):
			// 超时, 强制下线
			this.SendOne(user, fmt.Sprintf(ImConstant.IdleMsg+"\n"))
			user.Offline(this)
			return
		}
	}
}

func (this *Server) SendMessage(user *User, msg string) {
	if msg == ImConstant.WhoCommand {
		// 查询在线用户
		this.mapLock.Lock()
		for _, u := range this.OnlineMap {
			sendMsg := fmt.Sprintf("[%s]: %s online...\n", u.Addr, u.Name)
			this.SendOne(user, sendMsg)
		}
		this.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == ImConstant.RenameCommand {
		// 重命名当前用户
		// 格式：rename|姓名
		newName := strings.Split(msg, "|")[1]
		// 判断用户名是否存在
		_, ok := this.OnlineMap[newName]
		if ok {
			this.SendOne(user, fmt.Sprintf("user name [%s] has been used\n", newName))
		} else {
			this.mapLock.Lock()
			delete(this.OnlineMap, user.Name)
			this.OnlineMap[newName] = user
			user.Name = newName
			this.SendOne(user, fmt.Sprintf("user name has been changed to %s\n", newName))
			this.mapLock.Unlock()
		}
	} else if len(msg) > 4 && msg[:3] == ImConstant.ToCommand {
		// 私聊
		// 格式：to|name|msg
		// 1. 获取对象用户
		targetName := strings.Split(msg, "|")[1]
		if targetName == "" {
			this.SendOne(user, "invalid command, the correct format is \"to|name|msg\"\n")
		}
		targetUser, ok := this.OnlineMap[targetName]
		if !ok {
			this.SendOne(user, "target user does not exist\n")
			return
		}
		// 2. 发送消息
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendOne(user, "blank msg, skip\n")
		}
		sendMsg := fmt.Sprintf("receive message from [%s]: %s\n", user.Name, content)
		this.SendOne(targetUser, sendMsg)
	} else {
		// 广播
		sendMsg := fmt.Sprintf("[%s]: %s", user.Name, msg)
		this.BroadCast(sendMsg)
	}
}

func (this *Server) BroadCast(msg string) {
	this.Message <- msg
}

func (this *Server) SendOne(user *User, msg string) {
	_, _ = user.Conn.Write([]byte(msg))
}

// Start 启动服务器的接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	//close listen socket
	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	// listen message
	go this.MessageListener()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		//do handler
		go this.Handler(conn)
	}
}
