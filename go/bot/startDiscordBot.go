package bot

import (
	"Leetcode-or-Explode-Bot/db"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"google.golang.org/api/sheets/v4"
	"os"
	"os/signal"
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

	go dailyposts(sess)

	if err != nil {
		fmt.Println("Error creating Discord session,", err)
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
					Content: "You best start leetcoding soon... or else 😡",
				},
			})

		case "signup":
			fmt.Println("signup")

			lcUsername := i.ApplicationCommandData().Options[0].StringValue()

			var sheets *sheets.Service

			sheets, err = getGoogleSheets()
			if err != nil {
				fmt.Println("❌ Failed to initialize Google Sheets client:", err)
				return
			}
			createNewSheetWithTitle(sheets, spreadsheetID, i.Member.User.Username)

			if db.DoesExist(db.DB, "users", "user_id", lcUsername) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "This user already exists: " + lcUsername,
					},
				})
				return
			}

			err := db.AddUser(db.DB,
				lcUsername,
				false,
				0,
				"DEFAULT",
				i.Member.User.ID,
				i.GuildID,
				i.Member.User.Username)
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

			//

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "✅ Signed up as " + lcUsername,
				},
			})

		case "delete":
			fmt.Println("delete")
			db.DeleteUserByDiscordID(db.DB, i.Member.User.ID)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "User deleted! you can now re-signed up!",
				},
			})

		case "random leetcode":
			fmt.Println("Random Leecoset")

			var lcURL string // default is 0

			lcURL = db.GetRandomSkewedLeetcode(db.DB, i.Member.User.ID)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: lcURL,
				},
			})
		}

		// TODO: Add a lockup to leetcode site (or api) to see if this user exists or not
	})

	sess.Identify.Intents = discordgo.IntentsGuildMessages

	err = sess.Open()
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
		},
		{
			Name:        "random leetcode",
			Description: "Get a radome Leetcode with skewed probability in favour of least confident past Leetcodes",
		},
	}

	for _, cmd := range commands {
		_, err := sess.ApplicationCommandCreate(sess.State.User.ID, "1392352918425960509", cmd)
		if err != nil {
			fmt.Printf("❌ Cannot create slash command '%s': %v\n", cmd.Name, err)
		}
	}
	// --------------------- End of commands -----------------------

	err = sess.GuildMemberNickname("1392352918425960509", "@me", "Leetcode Warden 😭")
	if err != nil {
		fmt.Println("Error creating nickname", err)
	}

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
