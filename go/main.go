package main

import (
	"Leetcode-or-Explode-Bot/bot"
	"Leetcode-or-Explode-Bot/db"
	"database/sql"
)

var Conn *sql.DB

func main() {

	db.Init()

	go bot.StartDiscordBot() // requires wifie
	go bot.StartChromeAPIServer()

	select {} // cleaner than Sleep for long-running goroutines
}
