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

func SetPokerusLock(conn *sql.DB, time time.Time) error {
	if s, err := time.MarshalText(); err != nil {
		return err
	} else {
		_, err := conn.Exec("insert into global_properties (key,value) values ('POKERUS_LOCK', $1)", s)
		return err
	}
}

// Returns the value of the POKERUS LOCK.
// This lock should be set when Pokérus host is announced in a Discord guild, see SetPokerusLock.
// Returns the time stamp zero if this lock has never been set.
func GetPokerusLock(conn *sql.DB) (time.Time, error) {
	var time time.Time
	row := conn.QueryRow("select value from global_properties where key='POKERUS_LOCK'")
	var b []byte
	err := row.Scan(&b)
	if err != nil && err != sql.ErrNoRows {
		// an error occured
		return time, err
	} else if err == sql.ErrNoRows {
		// the lock was never set, return the zero-timestapm and no error
		return time, nil
	} else {
		// unmarshal the timestamp
		err = time.UnmarshalText(b)
		return time, err
	}
}

func RemoveUserByDiscordId(conn *sql.DB, discordId string) (bool, error) {
	result, err := conn.Exec("delete from usernames where discord = $1", discordId)
	if err != nil {
		// internal error
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// internal error
		return false, err
	}
	if rowsAffected == 0 {
		// user was not found in the database
		return false, nil
	}
	// user was deleted
	return true, nil
}
