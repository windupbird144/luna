package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"luna/operations"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
)

type LunaConfig struct {
	ApplicationId   string
	BotToken        string
	PostgresUri     string
	HugDirectory    string
	Server          string
	PokerusLockTime int
}

var config LunaConfig
var s *discordgo.Session
var db *sql.DB
var err error

func init() { rand.Seed(time.Now().UnixNano()) }

func init() { flag.Parse() }

func init() {
	possibleConfigFileLocations := []string{"/etc/luna/config.toml"}
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		possibleConfigFileLocations = append(possibleConfigFileLocations, os.ExpandEnv("$XDG_CONFIG_HOME/luna/config.toml"))
	} else {
		possibleConfigFileLocations = append(possibleConfigFileLocations, os.ExpandEnv("$HOME/.config/luna/config.toml"))
	}
	if os.Getenv("LUNA_CONFIG") != "" {
		possibleConfigFileLocations = append(possibleConfigFileLocations, os.ExpandEnv("$LUNA_CONFIG"))
	}
	log.Printf("checking %v locations: %v", len(possibleConfigFileLocations), possibleConfigFileLocations)
	var configLocation string
	for _, location := range possibleConfigFileLocations {
		if _, err := os.Stat(location); os.IsNotExist(err) {
			continue
		} else {
			configLocation = location
			break
		}
	}
	if configLocation == "" {
		log.Fatal("None of the config locations exist!")
	}
	if _, err = toml.DecodeFile(configLocation, &config); err != nil {
		log.Fatalf("could not decode the config file at %v: %v", configLocation, err)
	}
}

func init() {
	var err error
	s, err = discordgo.New("Bot " + config.BotToken)
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
		{
			Name:        "hyperbeam",
			Description: "fire off a hyperbeam",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "target",
					Description: "the target of your hyperbeam",
					Required:    true,
				},
			},
		},
		{
			Name:        "setreminder",
			Description: "set a reminder",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "what",
					Description: "what you want to be reminded of",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "when",
					Description: "e.g. in 1 hour",
					Required:    true,
				},
			},
		},
		{
			Name:        "forgetme",
			Description: "do not receive any more pokerus notifications",
			Options:     []*discordgo.ApplicationCommandOption{},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"hug": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Println("received a /hug command")
			id := i.ApplicationCommandData().Options[0].UserValue(nil).ID
			log.Printf("hugging user %s\n", id)
			gif, err := operations.ReadRandomFile(config.HugDirectory)
			if err != nil {
				log.Printf("error getting gif: %v\n", err)
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			_, err = s.FollowupMessageCreate(config.ApplicationId, i.Interaction, true, &discordgo.WebhookParams{
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
			username := operations.NewUsername(i.ApplicationCommandData().Options[0].StringValue())
			ok, err := operations.UserExists(username)
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
			reply := operations.CreateMapping(db, i.Member.User.ID, username)
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
		"hyperbeam": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if target := i.ApplicationCommandData().Options[0].UserValue(nil); target != nil {
				hyperbeam := operations.NewHyperBeam()
				reply := fmt.Sprintf("%s users Hyper Beam on %s! %s takes %d damage!", i.Member.Mention(), target.Mention(), target.Mention(), hyperbeam.ActualDamage())
				if hyperbeam.ActualDamage() >= 100 {
					reply = reply + fmt.Sprintf(" It's a critical hit! %s fainted!", target.Mention())
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: reply,
					},
				})
			}
		},
		"setreminder": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			opts := i.ApplicationCommandData().Options
			what := opts[0].StringValue()
			when := opts[1].StringValue()
			if what != "" && when != "" {
				if reminder, err := operations.NewReminder(i.Member.User.ID, when, what, time.Now()); err != nil {
					log.Printf("error creating reminder from options %v", err)
				} else {
					if err := operations.InsertReminder(db, i.GuildID, reminder); err != nil {
						log.Printf("error saving reminder %v", err)
					} else {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: fmt.Sprintf("got it, i will remind you on %s", reminder.Due.Format(time.RFC1123)),
							},
						})
					}
				}
			} else {
				log.Printf("error parsing arguments what='%v' when ='%v'", what, when)
			}
		},
		"forgetme": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			success, err := operations.DeleteFromRemindersTableByDiscordId(db, i.Member.User.ID)
			var msg string
			if err != nil {
				msg = "internal error :("
			} else if success {
				msg = "ok, you will not receive any more pokerus notifications!"
			} else {
				msg = "it appears notifications were never turned on for your discord user"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
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
	// Connect to Postgres, retry every 5 seconds until the connection succeeds, then continue
	for {
		db, err = sql.Open("postgres", config.PostgresUri)
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
	// To make Luna announce Pokérus, send a HTTP request to the port passed in the command line options
	startServer := func() {
		log.Println("Starting the Pokérus server")
		http.HandleFunc("/pokerus", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Received a request to announce Pokérus")
			// Get the Pokerus lock. If the lock was set less than ten minutes ago, do not announce pokerus.
			lock_time, err := operations.GetPokerusLock(db)
			if err != nil {
				log.Printf("Error getting the Pokerus lock: %v\n", err)
			} else {
				log.Printf("Got Pokerus lock: %v\n", lock_time)
			}
			duration := time.Now().Sub(lock_time)
			minutes_since_last_announcement := duration.Minutes()

			if minutes_since_last_announcement < float64(config.PokerusLockTime) {
				log.Printf("Refusing to announce Pokerus because Pokerus was already announced %v minutes ago.", minutes_since_last_announcement)
				return
			} else {
				log.Printf("Proceeding to announce - Minutes since last announcement: %v\n", minutes_since_last_announcement)
			}

			// get the current pokerus holder
			log.Printf("Fetching Pokerus host")
			pokerus, err := operations.Pokerus()
			if err != nil {
				log.Printf("Could not get the Pokerus host: %v\n", err)
				return
			} else {
				log.Printf("Got Pokerus host: %v\n", pokerus)
			}

			// get the discord ID of the host
			discordId := operations.GetDiscordId(db, operations.NewUsername(pokerus.Name))

			// iterate over guilds
			log.Printf("begin iterating over guilds")
			for _, guild := range s.State.Guilds {
				log.Printf("looking for a channel called 'rus-alert' in guild %v", guild.ID)
				// find the channel named rus-alert
				pokerusChannelFound := false
				for _, channel := range guild.Channels {
					if strings.Contains(channel.Name, "rus-alert") {
						pokerusChannelFound = true
						log.Printf("found channel in guild ID %v with channel ID %v", guild.ID, channel.ID)
						// check if the pokerus host is in the guild
						// they are in the guild if s.GuildMember returns a non-nil member and a nil error
						hostInGuild := false
						if discordId != "" {
							member, err := s.GuildMember(guild.ID, discordId)
							if err != nil {
								log.Printf("Error getting guild member: %v\n", err)
							}
							hostInGuild = member != nil && err == nil
						}
						// Construct the message (if the host is in the guild: with a ping)
						var messagePart string
						if hostInGuild {
							messagePart = fmt.Sprintf("<@%s>", discordId)
						} else {
							messagePart = pokerus.Name
						}
						message := fmt.Sprintf("%s has Pokérus <%s>", messagePart, pokerus.Url)
						log.Printf("sending pokerus message")
						if _, err := s.ChannelMessageSend(channel.ID, message); err == nil {
							log.Printf("Successfully sent pokerus message")
						} else {
							log.Printf("Could NOT send the pokerus message, error was %v", err)
						}
					}
				}
				if !pokerusChannelFound {
					log.Printf("warning: did not find any channel containing rus-alert in guild %v", guild.ID)
				}
			}

			// set the pokerus lock
			err = operations.SetPokerusLock(db, time.Now())
			if err != nil {
				log.Printf("Error setting the Pokerus lock: %v", err)
			} else {
				log.Printf("Successfully set the Pokerus lock")
			}
		})
		http.HandleFunc("/reminders", func(w http.ResponseWriter, r *http.Request) {
			// Check for due reminders
			tim := time.Now()
			if reminders, err := operations.GetDueReminders(db, tim); err != nil {
				log.Printf("error checking reminders %v", err)
			} else {
				log.Printf("Got %v due reminders", len(reminders))
				for _, r := range reminders {
					msg := fmt.Sprintf("<@%s> Here is your reminder: '%s'", r.Reminder.DiscordId, r.Reminder.Text)
					// This will become a problem if you add luna to more than 1 channel
					log.Printf("Trying to find guild with id %v", r.GuildID)
					if guild, err := s.State.Guild(r.GuildID); err != nil {
						log.Printf("error finding guild %v", err)
					} else {
						log.Printf("Found guild ID %v", r.GuildID)
						log.Printf("Looking for a channel called 'bot' in guild ID %v", r.GuildID)
						for _, channel := range guild.Channels {
							if channel.Name == "bot" {
								log.Printf("Sending reminder message")
								if _, err := s.ChannelMessageSend(channel.ID, msg); err != nil {
									log.Printf("error sending message %v", err)
								} else {
									if err := operations.DeleteDueReminders(db, tim); err != nil {
										log.Printf("error deleting reminders %v", err)
									} else {
										log.Printf("Successfully sent and deleted reminder")
									}
								}
							}
						}
					}
				}
			}
		})
		err = http.ListenAndServe(config.Server, nil)
	}
	go startServer()

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Removing previous commands...")
	registeredCommands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Adding commands...")
	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
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
