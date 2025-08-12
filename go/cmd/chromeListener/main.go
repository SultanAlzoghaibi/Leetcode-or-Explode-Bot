package main

import (
	"Leetcode-or-Explode-Bot/internal/chrome"
	"Leetcode-or-Explode-Bot/internal/db"
	"database/sql"
	"fmt"
	"runtime/debug"
	"time"
)

var Conn *sql.DB

func main() {
	fmt.Println("main chrome started")
	db.Init()
	go recoverer(5, 1, chrome.StartChromeAPIServer)
	select {} // cleaner than Sleep for long-running goroutines
}

func recoverer(maxPanics, id int, f func()) {
	// from stackoverflow
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("HERE", id)
			fmt.Println(err)
			fmt.Println("Stack trace:\n", string(debug.Stack()))
			fmt.Println(" number of cyles left", maxPanics)

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
