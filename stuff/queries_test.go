package stuff

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestAddReminder(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost user=luna dbname=luna sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}
	tim, err := time.Parse("2006", "2000")
	if err != nil {
		t.Fatal(err)
	}
	testCases := map[string]int{
		"1999": 0,
		"2001": 1,
	}
	// Insert one reminder at the year 2000
	if err = InsertReminder(db, "", Reminder{DiscordId: "1234", Due: tim, Text: "hello"}); err != nil {
		t.Fatal(err)
	} else {
		// In 1999 there should be no due reminders
		// In 2001 there should be one
		for year, expectedReminders := range testCases {
			if tim, err = time.Parse("2006", year); err != nil {
				t.Fatal(err)
			} else {
				if reminders, err := GetDueReminders(db, tim); err != nil {
					t.Fatal(err)
				} else {
					assert.Equal(t, expectedReminders, len(reminders))
				}
			}
		}
		// After deleteing due reminders in 1999, the one in 2000 should remain
		// After deleting due reminders in 2001, the one in 2000 should be gone
		if future, err := time.Parse("2006", "2000"); err != nil {
			t.Fatal(err)
		} else {
			for year, i := range testCases {
				expectedRemaining := 1 - i
				if tim, err = time.Parse("2006", year); err != nil {
					t.Fatal(err)
				} else {
					if err := DeleteDueReminders(db, tim); err != nil {
						t.Fatal(err)
					} else {
						if reminder, err := GetDueReminders(db, future); err != nil {
							t.Fatal(err)
						} else {
							assert.Equal(t, expectedRemaining, len(reminder))
						}
					}
				}
			}
		}
	}
}
