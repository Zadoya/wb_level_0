package main

import (
	"wb_level_0/internal/server"
)

func main() {
	server := server.NewServer()
	server.Start()
	server.Shutdown()
}