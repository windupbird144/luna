package operations

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test if the due date of the reminder is what you expect
func TestNewReminder(t *testing.T) {
	// We create a reminder on January 1 with the due date "in 48 hours"
	// The reminder should be due on January 3
	january1 := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	january3 := time.Date(2000, time.January, 3, 0, 0, 0, 0, time.UTC)

	reminder, _ := NewReminder("", "in 48 hours", "", january1)
	assert.Equal(t, january3, reminder.Due, "wrong due date in NewReminder")
}

func TestNewReminderBadDueDate(t *testing.T) {
	reminder, err := NewReminder("", "in ??? minutes", "", time.Now())
	assert.Error(t, err, "expected an error because of a bad timestamp, found %v", reminder)
}
