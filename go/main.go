package main

import (
	"time"
)

func main() {
	go startDiscordBot()
	go startChromeAPIServer()
	time.Sleep(100 * time.Second) // Let goroutines print something

}
