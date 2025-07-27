package main

import (
	"Leetcode-or-Explode-Bot/bot"
	"Leetcode-or-Explode-Bot/db"
	"database/sql"
	"fmt"
	"runtime/debug"
)

var Conn *sql.DB

func main() {
	fmt.Println("main started")
	db.Init()
	go recoverer(100, 1, bot.StartChromeAPIServer)
	//go recoverer(100, 2, bot.StartDiscordBot)

	select {} // cleaner than Sleep for long-running goroutines
}

func recoverer(maxPanics, id int, f func()) {
	// from stackoverflow
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("HERE", id)
			fmt.Println(err)
			fmt.Println("Stack trace:\n", string(debug.Stack()))

			if maxPanics == 0 {
				panic("TOO MANY PANICS")
			} else {
				go recoverer(maxPanics-1, id, f) // restart goroutine with 1 fewer retry
			}
		}
	}()
	f()
}
