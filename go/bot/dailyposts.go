package bot

import (
	"Leetcode-or-Explode-Bot/db"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

func dailyposts(s *discordgo.Session) {
	fmt.Println("✅ Daily post loop started")

	sleep := 24 * time.Second

	for {
		now := time.Now()
		date := now.Format("2006-01-02")

		// 1. Get stats
		res := db.GetAllDailyLeets(db.DB, date)
		fmt.Println(res)

		// 2. Send daily stats message to a channel
		s.ChannelMessageSend("1395556314951974972", "res")

		// 3. Send leaderboard
		leaderboardmsg := DisplayLeaderboard(db.GetLeaderboard(db.DB))
		s.ChannelMessageSend("1395556365623234600", leaderboardmsg)

		time.Sleep(sleep)
	}
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
