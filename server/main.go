package main

import (
	"server/core"
)

const (
	IP   = "127.0.0.1"
	PORT = 80
)

func main() {
	server := core.Server{}

	server.ConfigServer()
	server.InitServer()
	server.StartServer()
}
