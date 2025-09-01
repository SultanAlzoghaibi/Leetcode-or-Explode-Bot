package bot

import (
	db2 "Leetcode-or-Explode-Bot/internal/db"
	"Leetcode-or-Explode-Bot/internal/shared"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"sort"
	"strings"
	"syscall"
	"time"
)

const prefix = "lcc"

type Command struct {
}

func StartDiscordBot() {
	godotenv.Load()
	discToken := os.Getenv("DISCORD_TOKEN")

	fmt.Println("token:" + discToken)
	sess, err := discordgo.New("Bot " + discToken) // I think this turn on the bot
	if err != nil {
		log.Printf("discordgo.New failed: %v", err)
		return
	}
	sess.Client = &http.Client{Timeout: 30 * time.Second}

	go dailyposts(sess) // temp oof  undo
	//TODO: UNDO

	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		fmt.Println("Stack trace:\n", string(debug.Stack()))
	}
	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		//todo: verify that this thread depedan ton the seshpoing being active
		/*
			userID := i.Member.User.ID
			guildID := i.GuildID
			channelID := i.ChannelID
			lcUsername := i.ApplicationCommandData().Options[0].StringValue()

			m holds data like:
				•	m.Content → the actual message text (string)
				•	m.Author.ID → who sent it
				•	m.ChannelID → where it was sent
				•	m.GuildID
		*/

		switch i.ApplicationCommandData().Name {

		case "ping":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "no pong here, ONLY the Leetcode Warden, You best start leetcoding soon... or else! 😡",
				},
			})

		case "signup":
			fmt.Println("signup")

			lcUsername := i.ApplicationCommandData().Options[0].StringValue()

			if db2.DoesExist(db2.DB, "users", "user_id", lcUsername) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "This user already exists: " + lcUsername,
					},
				})
				return
			}

			var sheets *sheets.Service

			sheets, err = shared.GetGoogleSheets()
			if err != nil {
				fmt.Println("❌ Failed to initialize Google Sheets client:", err)
				return
			}

			go shared.CreateNewSheetWithTitle(sheets, shared.SpreadsheetID, i.Member.User.Username)

			err := db2.AddUser(db2.DB,
				lcUsername,
				false,
				0,
				"DEFAULT",
				i.Member.User.ID,
				i.GuildID,
				i.Member.User.Username,
				0)
			if err != nil {
				if strings.Contains(err.Error(), "Error 1062") {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "⚠️ You've already signed up!",
							Flags:   1 << 6, // ephemeral message
						},
					})
					return
				}
				fmt.Println("❌ Failed to add user:", err)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "✅ Signed up as " + lcUsername,
				},
			})

		case "delete":
			fmt.Println("delete")
			confirmInput := i.ApplicationCommandData().Options[0].StringValue()
			if strings.ToLower(confirmInput) != "confirm" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "⚠️ You must type `confirm` to delete your data. This action is irreversible.",
						Flags:   1 << 6,
					},
				})
				return
			}

			// ✅ Check if user exists before proceeding
			exists := db2.DoesExist(db2.DB, "users", "discord_user_id", i.Member.User.ID)
			if !exists {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "⚠️ No record found for this user. Have you signed up?",
						Flags:   1 << 6,
					},
				})
				return
			}

			db2.DeleteUserByDiscordID(db2.DB, i.Member.User.ID)
			var sheets *sheets.Service

			sheets, err = shared.GetGoogleSheets()
			if err != nil {
				fmt.Println("❌ Failed to initialize Google Sheets client:", err)
				return
			}
			err := shared.DeleteSheetByTitle(sheets, shared.SpreadsheetID, i.Member.User.Username)
			if err != nil {
				log.Fatal(err)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "User deleted! you can now re-signed up!",
				},
			})

		case "random-leetcode":
			fmt.Println("Random-Leecoset")

			userID, err := db2.GetUserIDwithDiscordID(db2.DB, i.Member.User.ID)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "❌ Could not find your user record. Make sure you signed up.",
						Flags:   1 << 6,
					},
				})
				return
			}

			lcURL := db2.GetRandomSkewedLeetcode(db2.DB, userID)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: lcURL,
				},
			})

		case "status":

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "comming SOON",
				},
			})

		case "stalk":
			loc, _ := time.LoadLocation("America/Los_Angeles")
			date := time.Now().In(loc).Format("2006-01-02")

			dailyInfoMap, err := db2.GetLCFromAllUsersToday(db2.DB, date)
			if err != nil {
				fmt.Println(err)
				return
			}

			report := BuildDailyReport(dailyInfoMap)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: report,
					Flags:   1 << 6, // ephemeral message
				},
			})

		}

		// TODO: Add a lockup to leetcode site (or api) to see if this user exists or not
	})

	sess.Identify.Intents = discordgo.IntentsGuildMessages

	err = sess.Open()
	if err != nil {
		fmt.Println("Failed to open session:", err)
		fmt.Println("Stack trace:\n", string(debug.Stack()))
		return
	}
	// Register slash commands -------------------------------
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Replies with pong",
		},
		{
			Name:        "signup",
			Description: "Sign up with your LeetCode username",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "userid",
					Description: "Your LeetCode user name",
					Required:    true,
				},
			},
		},
		{
			Name:        "delete",
			Description: "Delete your self from the database",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "confirm",
					Description: "⚠️ Type `confirm` to delete all your data and spreadsheet",
					Required:    true,
				},
			},
		},
		{
			Name:        "random-leetcode",
			Description: "Get a radome Leetcode with skewed probability in favour of least confident past Leetcodes",
		},
		{
			Name:        "stalk",
			Description: "Get a list of everyone's solved Leetcode questions solved today",
		},
		//{
		//	Name:        "status",
		//	Description: "Change your ping status (DEFAULT or NO-PING)",
		//	Options: []*discordgo.ApplicationCommandOption{
		//		{
		//			Type:        discordgo.ApplicationCommandOptionString,
		//			Name:        "value",
		//			Description: "Choose your ping status (coming soon)",
		//			Required:    true,
		//			Choices: []*discordgo.ApplicationCommandOptionChoice{
		//				{
		//					Name:  "DEFAULT",
		//					Value: "DEFAULT",
		//				},
		//				{
		//					Name:  "NO-PING",
		//					Value: "NO-PING",
		//				},
		//			},
		//		},
		//	},
		//},
	}
	if os.Getenv("ENV") != "production" {
		fmt.Println("Loading environment variables from .env file")
		if err := godotenv.Load(".env"); err != nil {
			log.Println("❌ Failed to load .env:", err)
		} else {
			fmt.Println("✅ .env file loaded")
		}
	}

	var guildID string
	if os.Getenv("ENV") == "production" {
		guildID = "1377097284633886771"
	} else {
		guildID = "1392352918425960509"
	}

	for _, cmd := range commands {
		_, err := sess.ApplicationCommandCreate(sess.State.User.ID, guildID, cmd)
		if err != nil {
			fmt.Printf("❌ Cannot create slash command '%s' in guild %s: %v\n", cmd.Name, guildID, err)
		}
	}
	err = sess.GuildMemberNickname(guildID, "@me", "Leetcode Warden 😭")
	if err != nil {
		fmt.Printf("❌ Failed to set nickname in guild %s: %v\n", guildID, err)
	}

	// --------------------- End of commands -----------------------

	if err != nil {
		fmt.Println("❌ Failed to update nickname:", err)
	}
	if err != nil {
		fmt.Println("Error opening connection,", err)
	}
	defer sess.Close()
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

// SendDirectMessage sends a DM to a specific Discord user by ID.
func SendDirectMessage(s *discordgo.Session, userID, message string) error {
	// Create (or fetch) a DM channel with the user
	channel, err := s.UserChannelCreate(userID)
	if err != nil {
		return fmt.Errorf("failed to create DM channel with user %s: %w", userID, err)
	}

	// Send the message into the DM channel
	_, err = s.ChannelMessageSend(channel.ID, message)
	if err != nil {
		return fmt.Errorf("failed to send DM to user %s: %w", userID, err)
	}

	return nil
}

//todo: fix bug on first submission being error every time

// TODO: Add a stat strater with telemary (aksii art wort comes to worse)

func resetStreak(db *sql.DB, userID string) {
	updateQuery := `UPDATE users SET streak = 0 WHERE user_id = ?`
	stmt, err := db.Prepare(updateQuery)
	if err != nil {
		fmt.Println("Error resetting streak:", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(userID)
	if err != nil {
		log.Println("Error resetting streak:", err)
	}

}

func difficultyToEmoji(diff string) string {
	switch diff {
	case "EASY":
		return "🟩" // green square
	case "MEDIUM":
		return "🟨" // yellow square
	case "HARD":
		return "🟥" // red square
	default:
		return "⬜" // fallback
	}
}

func BuildDailyReport(dailyInfoMap map[string][]string) string {
	var sb strings.Builder
	if len(dailyInfoMap) == 0 {
		return "⚠️ No one solved any LeetCode problems today. BE THE FIRST 🥇"
	}

	// Step 1: build a slice of usernames with counts
	type userCount struct {
		username string
		count    int
	}
	var counts []userCount
	for u, entries := range dailyInfoMap {
		counts = append(counts, userCount{username: u, count: len(entries)})
	}

	// Step 2: sort by descending count
	sort.Slice(counts, func(i, j int) bool {
		if counts[i].count == counts[j].count {
			return counts[i].username < counts[j].username // tie-breaker: alphabetical
		}
		return counts[i].count > counts[j].count
	})

	// Step 3: render the report
	for _, uc := range counts {
		entries := dailyInfoMap[uc.username]
		sb.WriteString(fmt.Sprintf("**%s** (%d solved):\n", uc.username, uc.count))
		for _, entry := range entries {
			parts := strings.SplitN(entry, ":", 2)
			if len(parts) == 2 {
				problem := strings.TrimSpace(parts[0])
				diff := strings.TrimSpace(parts[1])
				emoji := difficultyToEmoji(strings.ToUpper(diff))
				sb.WriteString(fmt.Sprintf("   %s %s\n", problem, emoji))
			} else {
				sb.WriteString("   " + entry + "\n")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
