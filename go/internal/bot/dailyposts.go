package bot

import (
	db2 "Leetcode-or-Explode-Bot/internal/db"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
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
		var sleepDuration time.Duration
		if os.Getenv("ENV") == "test" {
			sleepDuration = 24 * time.Second
		} else {
			sleepDuration = time.Until(nextRun)
		}
		fmt.Printf("😴 Sleeping %v until 11:59 PM daily run at %v\n", sleepDuration, nextRun)
		time.Sleep(sleepDuration)

		// ----------- Now that it’s 11:59 PM, recalculate date -----------
		runDate := nextRun.In(loc)
		date := runDate.Format("2006-01-02")
		fmt.Printf("📅 Date: %s\n", date)

		// ----------- Do daily stuff -----------
		dailyStats := db2.GetAllDailyLeets(db2.DB, date)
		var channelDailyLCID, channelLeaderboardID string

		if os.Getenv("ENV") == "test" {
			channelDailyLCID = "1395556314951974972"     // test daily LC
			channelLeaderboardID = "1395556365623234600" // test leaderboard
		} else {
			channelDailyLCID = "1399588861461659678"     // prod daily LC
			channelLeaderboardID = "1399588897595588638" // prod leaderboard
		}
		fmt.Println("channels: ", channelDailyLCID, channelLeaderboardID)
		s.ChannelMessageSend(channelDailyLCID, DisplayDailylc(dailyStats))
		s.ChannelMessageSend(channelLeaderboardID, DisplayLeaderboard(db2.GetLeaderboard(db2.DB)))

		// Reset monthly LC if month changed
		if now.Add(1*time.Hour).Day() == 1 {
			db2.ResetMoLCA(db2.DB)
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

func DisplayLeaderboard(leaderboard []db2.LeaderEntry) string {
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

func DisplayDailylc(stats []db2.DailyStat) string {
	var res strings.Builder
	loc, _ := time.LoadLocation("America/Los_Angeles")
	now := time.Now().In(loc)
	res.WriteString(fmt.Sprintf("📅 Day %d — Daily Leetcode Records: \n\n", now.Day()))

	for _, stat := range stats {
		total := stat.Easy + stat.Medium + stat.Hard

		if total == 0 {
			resetStreak(db2.DB, stat.UserID)
		} else {
			db2.IncrementStreak(db2.DB, stat.UserID)
			stat.Streak++
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

			if stat.Streak >= 4 {
				res.WriteString(fmt.Sprintf(" |  Streak  %d 🔥", stat.Streak))
			}

			res.WriteString("\n")

		}

	}

	return res.String()
}
