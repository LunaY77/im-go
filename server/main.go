package main

import ImConstant "github.com/LunaY77/im-go/const"

func main() {
	server := NewServer(ImConstant.ServerIp, ImConstant.ServerPort)
	server.Start()
}
