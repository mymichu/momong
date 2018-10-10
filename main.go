package main

import (
	"local/momong/server"
	"runtime"
)

func main() {
	go server.Start()
	go server.StartTodo()
	runtime.Goexit()
}
