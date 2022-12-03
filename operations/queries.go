package operations

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Add a mapping from discord ID to username
func CreateMapping(conn *sql.DB, discordId string, username Username) string {
	conflict := ""

	// Check if user already has a pfq name
	conn.QueryRow("select username from usernames where discord = $1", discordId).Scan(&conflict)
	if conflict != "" {
		// update it
		_, err := conn.Exec("update usernames set username = $1 where discord = $2", string(username), discordId)
		if err != nil {
			return ""
		}
		return fmt.Sprintf("okay, i updated your username to %s! (previously %s)", string(username), conflict)
	}

	// Check for conflict again
	conn.QueryRow("select * from usernames where username = $1", discordId).Scan(&conflict)
	if conflict != "" {
		return "nope, i already have a discord user saved for this name!"
	}

	// Save user
	_, err := conn.Exec("insert into usernames (discord, username) values ($1, $2)", discordId, string(username))
	if err != nil {
		// Something went wrong!
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("got it, i saved you as %s", string(username))
}

// Get the discord ID for this PFQ username. Returns "" if no entry is found
func GetDiscordId(conn *sql.DB, username Username) string {
	discord := ""
	conn.QueryRow("select discord from usernames where username = $1", string(username)).Scan(&discord)
	return discord
}

func InsertReminder(conn *sql.DB, guildId string, reminder Reminder) error {
	_, err := conn.Exec("insert into reminders (guild_id, discord_id, due, text) values ($1, $2, $3, $4)", guildId, reminder.DiscordId, reminder.Due, reminder.Text)
	return err
}

type ReminderWithGuildId struct {
	GuildID  string
	Reminder Reminder
}

func GetDueReminders(conn *sql.DB, time time.Time) ([]ReminderWithGuildId, error) {
	results := make([]ReminderWithGuildId, 0)
	rows, err := conn.Query("select guild_id, discord_id, text from reminders where due < $1", time)
	if err != nil {
		return results, err
	}
	for rows.Next() {
		tmp := ReminderWithGuildId{}
		rows.Scan(&tmp.GuildID, &tmp.Reminder.DiscordId, &tmp.Reminder.Text)
		results = append(results, tmp)
	}
	return results, nil
}

func DeleteDueReminders(conn *sql.DB, time time.Time) error {
	_, err := conn.Exec("delete from reminders where due < $1", time)
	return err
}
