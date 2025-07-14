package main

import (
	"Leetcode-or-Explode-Bot/bot"
	"time"
)

func main() {
	go bot.startDiscordBot()
	go bot.startChromeAPIServer()
	time.Sleep(100 * time.Second) // Let goroutines print something

}
