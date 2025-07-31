package bot

import (
	"Leetcode-or-Explode-Bot/db"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
)

const prefix = "lcc"

type Command struct {
}

func StartDiscordBot() {
	godotenv.Load()
	discToken := os.Getenv("DISCORD_TOKEN")

	fmt.Println("token:" + discToken)
	sess, err := discordgo.New("Bot " + discToken) // I think this turn on the bot

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

			if db.DoesExist(db.DB, "users", "user_id", lcUsername) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "This user already exists: " + lcUsername,
					},
				})
				return
			}

			var sheets *sheets.Service

			sheets, err = getGoogleSheets()
			if err != nil {
				fmt.Println("❌ Failed to initialize Google Sheets client:", err)
				return
			}
			go createNewSheetWithTitle(sheets, spreadsheetID, i.Member.User.Username)

			err := db.AddUser(db.DB,
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
			exists := db.DoesExist(db.DB, "users", "discord_user_id", i.Member.User.ID)
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

			db.DeleteUserByDiscordID(db.DB, i.Member.User.ID)
			var sheets *sheets.Service

			sheets, err = getGoogleSheets()
			if err != nil {
				fmt.Println("❌ Failed to initialize Google Sheets client:", err)
				return
			}
			err := deleteSheetByTitle(sheets, spreadsheetID, i.Member.User.Username)
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

			userID, err := db.GetUserIDwithDiscordID(db.DB, i.Member.User.ID)
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

			lcURL := db.GetRandomSkewedLeetcode(db.DB, userID)

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
	serverIDs := []string{
		"1392352918425960509", // old server
		"1377097284633886771", // new server
	}

	for _, guildID := range serverIDs {
		for _, cmd := range commands {
			_, err := sess.ApplicationCommandCreate(sess.State.User.ID, guildID, cmd)

			if err != nil {
				fmt.Printf("❌ Cannot create slash command '%s' in guild %s: %v\n", cmd.Name, guildID, err)
			}
		}

		err := sess.GuildMemberNickname(guildID, "@me", "Leetcode Warden 😭")
		if err != nil {
			fmt.Printf("❌ Failed to set nickname in guild %s: %v\n", guildID, err)
		}
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
