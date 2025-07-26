package bot

import (
	"Leetcode-or-Explode-Bot/db"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

func dailyposts(s *discordgo.Session) {
	fmt.Println("✅ Daily post loop started")

	sleep := 12 * time.Second

	for {
		now := time.Now()
		date := now.Format("2006-01-02")

		// 1. Get stats
		dailyStats := db.GetAllDailyLeets(db.DB, date)
		fmt.Println(dailyStats)

		// 2. Send daily stats message to a channel
		s.ChannelMessageSend("1395556314951974972", DisplayDailylc(dailyStats))
		//TODO: change these to being the new discord Serverg

		// 3. Send leaderboard
		leaderboardmsg := DisplayLeaderboard(db.GetLeaderboard(db.DB))
		s.ChannelMessageSend("1395556365623234600", leaderboardmsg)

		currentMonth := now.Month()

		time.Sleep(sleep) // 12 hours
		//wasInative(db.DB, db.QueryAllSuerActivity(db.DB), s) //TODO commin soon feature of reminders
		time.Sleep(sleep) // 12 hours

		lastRecordedMonth := now.Month()

		if currentMonth != lastRecordedMonth {
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
