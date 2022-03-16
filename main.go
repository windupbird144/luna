package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"luna/stuff"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
)

// Bot parameters
var (
	ApplicationId = flag.String("app", "", "Discord application ID")
	BotToken      = flag.String("token", "", "Discord access token")
	PostgresUri   = flag.String("db", "", "Postgres URI")
	GuildID       = flag.String("guild", "", "(optional) guild ID for testing")
	HugDirectory  = flag.String("hugdir", "", "Directory containing hug gifs")
)

var s *discordgo.Session
var db *sql.DB
var err error

func init() { rand.Seed(time.Now().UnixNano()) }

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "hug",
			Description: "hug a friend",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "receiver",
					Description: "person who receives the hug",
					Required:    true,
				},
			},
		},
		{
			Name:        "pokerus",
			Description: "get a ping when you have pokerus",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "your pfq username",
					Required:    true,
				},
			},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"hug": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Println("received a /hug command")
			id := i.ApplicationCommandData().Options[0].UserValue(nil).ID
			log.Printf("hugging user %s\n", id)
			gif, err := stuff.ReadRandomFile(*HugDirectory)
			if err != nil {
				log.Printf("error getting gif: %v\n", err)
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			_, err = s.FollowupMessageCreate(*ApplicationId, i.Interaction, true, &discordgo.WebhookParams{
				Content: fmt.Sprintf("\\*hugs <@%s>\\*", id),
				Files: []*discordgo.File{
					{
						Name:        "hug.gif",
						ContentType: "image/gif",
						Reader:      gif,
					},
				},
			})
			if err != nil {
				log.Println(err)
			}
		},
		"pokerus": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			username := stuff.NewUsername(i.ApplicationCommandData().Options[0].StringValue())
			ok, err := stuff.UserExists(username)
			if err != nil {
				return
			}
			if !ok {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "sorry, i can't find anyone on pfq with that username.",
					},
				})
			}
			reply := stuff.CreateMapping(db, i.Member.User.ID, username)
			// Something went wrong
			if reply == "" {
				reply = "sorry, internal error :("
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: reply,
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	for {
		db, err = sql.Open("postgres", *PostgresUri)
		err = db.Ping()
		if err != nil {
			log.Printf("Postgres not connected, retry in 5 seconds %v", err)
			time.Sleep(5 * time.Second)
		} else {
			log.Println("Postgres connected")
			break
		}
	}

	// Announce Pokerus
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		ch := make(chan stuff.User)
		go stuff.PokeursChannel(ch)
		for {
			pokerus := <-ch
			// find all pokerus channels
			pokerusChannels := make([]string, 0)
			for _, guild := range s.State.Guilds {
				for _, channel := range guild.Channels {
					if channel.Name == "rus-alert" {
						pokerusChannels = append(pokerusChannels, channel.ID)
					}
				}
			}
			// get the member ID if available
			discordId := stuff.GetDiscordId(db, stuff.NewUsername(pokerus.Name))
			var message string
			if discordId == "" {
				message = pokerus.Name
			} else {
				message = fmt.Sprintf("<@%s>", discordId)
			}
			message = fmt.Sprintf("%s has Pokérus <%s>", message, pokerus.Url)
			// announce pokerus in every channel
			for _, channelId := range pokerusChannels {
				s.ChannelMessageSend(channelId, message)
			}
		}
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Removing previous commands...")
	registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Adding commands...")
	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutdowning")
}
