package bot

import (
	"Leetcode-or-Explode-Bot/db"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
	_ "time/tzdata"
)

func dailyposts(s *discordgo.Session) {

	fmt.Println("✅ Daily post loop started")

	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatalf("❌ Failed to load timezone: %v", err)

	}

	for {
		// ----------- Sleep until 11:59 PM -----------
		now := time.Now().In(loc)
		nextRun := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 0, 0, loc)
		if now.After(nextRun) {
			// If it's already past 11:59 PM today, schedule for tomorrow
			nextRun = nextRun.Add(24 * time.Hour)
		}
		sleepDuration := time.Until(nextRun)
		fmt.Printf("😴 Sleeping %v until 11:59 PM daily run at %v\n", sleepDuration, nextRun)
		time.Sleep(sleepDuration)

		// ----------- Now that it’s 11:59 PM, recalculate date -----------
		runDate := nextRun.In(loc)
		date := runDate.Format("2006-01-02")
		fmt.Printf("📅 Date: %s\n", date)

		// ----------- Do daily stuff -----------
		dailyStats := db.GetAllDailyLeets(db.DB, date)
		fmt.Println("dailyStats: ", dailyStats)

		s.ChannelMessageSend("1399588861461659678", DisplayDailylc(dailyStats))
		s.ChannelMessageSend("1399588897595588638", DisplayLeaderboard(db.GetLeaderboard(db.DB)))

		// Reset monthly LC if month changed
		if now.Add(1*time.Hour).Day() == 1 {
			db.ResetMoLCA(db.DB)
		}
	}
}

// setup a streak system ties to roles to reward stresks

func wasInative(db *sql.DB, hashmap map[string]bool, s *discordgo.Session) {
	var sb strings.Builder
	sb.WriteString("📢 **Ping of Shame** — No Leetcodes in the last 4 days!\n\n")

	for userID := range hashmap {
		sb.WriteString(fmt.Sprintf("<@%s> ", userID))
	}

	s.ChannelMessageSend("1395556365623234600", sb.String())
}

type LeaderEntry struct {
	Username   string
	MoLCAmount uint8
}

func DisplayLeaderboard(leaderboard []db.LeaderEntry) string {
	var res strings.Builder
	emojis := []string{"🥇", "🥈", "🥉"}

	res.WriteString("📊 Daily Leaderboard:\n")
	for i, entry := range leaderboard {
		var rank string
		if i < len(emojis) {
			rank = emojis[i]
		} else {
			rank = fmt.Sprintf("%d.", i+1)
		}
		res.WriteString(fmt.Sprintf("%s %s — %d\n", rank, entry.Username, entry.MoLCAmount))
	}

	return res.String()
}

func DisplayDailylc(stats []db.DailyStat) string {
	var res strings.Builder
	res.WriteString(fmt.Sprintf("📅 Day %d — Daily Leetcode Records: \n\n", time.Now().Day()))

	for _, stat := range stats {

		total := stat.Easy + stat.Medium + stat.Hard

		if total == 0 {
			resetStreak(db.DB, stat.UserID)
		} else {
			db.IncrementStreak(db.DB, stat.UserID)
		}

		res.WriteString(fmt.Sprintf(" %s — **%d** today | **%d** this month:\n", stat.Username, total, stat.MonthlyLC))
		if stat.Easy > 0 {
			res.WriteString(fmt.Sprintf("  🟩: %d", stat.Easy))
		}
		if stat.Medium > 0 {
			res.WriteString(fmt.Sprintf("  🟨: %d ", stat.Medium))
		}
		if stat.Hard > 0 {
			res.WriteString(fmt.Sprintf("  🟥: %d ", stat.Hard))
		}

		if stat.Streak > 3 {
			res.WriteString(fmt.Sprintf(" |  Streak  %d 🔥", stat.Streak))
		}

		res.WriteString("\n")

	}

	return res.String()
}
