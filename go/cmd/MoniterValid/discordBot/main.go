package main

import (
	bot2 "Leetcode-or-Explode-Bot/internal/bot"
	"Leetcode-or-Explode-Bot/internal/db"
	"database/sql"
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

var Conn *sql.DB

func main() {
	fmt.Println("main bot star")
	db.Init()

	recoverer(3, 2, bot2.StartDiscordBot)

	// We crahing the pot'
	fmt.Println("\n Too many panics, exiting...")

	time.Sleep(2 * time.Second)
	os.Exit(1)
}

func recoverer(maxPanics, id int, f func()) {
	// from stackoverflow
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("HERE", id)
			fmt.Println("number of panicked: ", maxPanics)
			fmt.Println(err)
			fmt.Println("Stack trace:\n", string(debug.Stack()))
			loc, _ := time.LoadLocation("America/Los_Angeles")
			fmt.Println("the time is", time.Now().In(loc).Format("2006-01-02 15:04:05 MST"))

			if maxPanics == 0 {
				panic("TOO MANY PANICS")
			} else {
				go recoverer(maxPanics-1, id, f) // restart goroutine with 1 fewer retry
			}
		}
	}()
	f()
}
