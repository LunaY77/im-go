package main

import (
	"fmt"
	"io"
	"net"

	ImConstant "github.com/LunaY77/IM-go/const"
)

type User struct {
	Name      string
	Addr      string
	Chan      chan string
	Conn      net.Conn
	HeartBeat chan bool
}

func NewUser(conn net.Conn, heartBeat chan bool) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:      userAddr,
		Addr:      userAddr,
		Chan:      make(chan string),
		Conn:      conn,
		HeartBeat: heartBeat,
	}
	go user.ListenMessage()
	return user
}

func (this *User) ListenMessage() {
	for {
		msg, ok := <-this.Chan
		if !ok {
			return
		}
		_, _ = this.Conn.Write([]byte(msg + "\n"))
	}
}

// Online 用户上线
func (this *User) Online(server *Server) {
	// 用户上线
	server.mapLock.Lock()
	server.OnlineMap[this.Name] = this
	server.mapLock.Unlock()

	// 广播消息
	server.SendMessage(this, "online")
}

// Offline 用户下线
func (this *User) Offline(server *Server) {
	// 删除用户
	server.mapLock.Lock()
	if _, ok := server.OnlineMap[this.Name]; ok {
		delete(server.OnlineMap, this.Name)
	}
	server.mapLock.Unlock()

	server.SendMessage(this, "offline")
	_ = this.Conn.Close()
}

// Message 用户发送消息
func (this *User) Message(server *Server) {
	buf := make([]byte, ImConstant.BufLen)
	for {
		n, err := this.Conn.Read(buf)
		if n == 0 {
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("Conn Read error: ", err)
			return
		}
		// 获取用户消息
		msg := string(buf[:n-1])
		// 广播消息
		server.SendMessage(this, msg)
		// 心跳
		this.Live()
	}
}

func (this *User) Live() {
	this.HeartBeat <- true
}
