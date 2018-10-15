package main

import (
	"local/momong/server"
	"local/momong/server_todo"
	"runtime"
)

func main() {
	go server.Start()
	go server_todo.Start()
	runtime.Goexit()
}
