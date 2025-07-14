package main

import (
	"Leetcode-or-Explode-Bot/bot"
	"time"
)

func main() {
	go bot.StartDiscordBot()
	go bot.StartChromeAPIServer()
	time.Sleep(100 * time.Second) // Let goroutines print something

}
