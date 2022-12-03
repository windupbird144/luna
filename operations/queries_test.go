package operations

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// Test that reminders are added to the database correctly
func TestInsertReminder(t *testing.T) {
	// test fixtures
	db, _ := sql.Open("postgres", "host=localhost user=postgres port=24019 sslmode=disable")
	defer db.Close()
	guild := "testGuild"
	discordId := "testId"
	due := time.Now()

	// delete old reminders and add one
	db.Exec("delete from reminders;")
	InsertReminder(db, guild, Reminder{
		DiscordId: discordId,
		Due:       due,
		Text:      "",
	})

	// assert that there is 1 reminder in the databse
	got := 0
	db.QueryRow("select count(*) from reminders").Scan(&got)
	assert.Equal(t, 1, got, "wrong number of reminders in the database")
}

// Test that the # of due reminders are correct
func TetsGetDueReminders(t *testing.T) {
	// test fixtures
	db, _ := sql.Open("postgres", "host=localhost user=postgres port=24019 sslmode=disable")
	defer db.Close()
	guild := "testGuild"
	discordId := "testId"

	timeStamp2000 := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	timeStamp2001 := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	timeStamp2002 := time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC)
	timeStamp2003 := time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC)

	// we clear the table and add two reminders: one due in 2000 and another in 2002.
	// in the year 2001 (between the two) one reminder should be due.
	// in the year 2003 (after the two) both should be due.
	db.Exec("delete from reminders;")
	InsertReminder(db, guild, Reminder{
		DiscordId: discordId,
		Due:       timeStamp2000,
		Text:      "",
	})
	InsertReminder(db, guild, Reminder{
		DiscordId: discordId,
		Due:       timeStamp2002,
		Text:      "",
	})

	// in 2001 there should be 1 due reminder
	reminders, _ := GetDueReminders(db, timeStamp2001)
	got := len(reminders)
	assert.Equal(t, 1, got, "wrong number of due reminders in 2001")

	// in 2003 there should be 2 due reminders
	reminders, _ = GetDueReminders(db, timeStamp2003)
	got = len(reminders)
	assert.Equal(t, 2, got, "wrong number of due reminders in 2003")
}

func TestDeleteDueReminders(t *testing.T) {
	db, _ := sql.Open("postgres", "host=localhost user=postgres port=24019 sslmode=disable")
	defer db.Close()
	guild := "testGuild"
	discordId := "testId"

	timeStamp2000 := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	timeStamp2001 := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	timeStamp2002 := time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC)

	// we clear the table and add two reminders: one due in 2000 and another in 2002.
	// in the year 2001 (between the two) one reminder should be due.
	// in the year 2003 (after the two) both should be due.
	db.Exec("delete from reminders;")
	InsertReminder(db, guild, Reminder{
		DiscordId: discordId,
		Due:       timeStamp2000,
		Text:      "",
	})
	InsertReminder(db, guild, Reminder{
		DiscordId: discordId,
		Due:       timeStamp2002,
		Text:      "",
	})

	// in 2001 one reminder should be removed and the other should be left
	DeleteDueReminders(db, timeStamp2001)
	got := 0
	db.QueryRow("select count(*) from reminders").Scan(&got)
	assert.Equal(t, 1, got, "wrong number of due reminders after deleteing due reminders")
}
