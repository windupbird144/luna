package stuff

import (
	"time"

	"github.com/tj/go-naturaldate"
)

type Reminder struct {
	DiscordId string
	Due       time.Time
	Text      string
}

func NewReminder(discordId string, due string, text string, currentTime time.Time) (Reminder, error) {
	// The default direction is future unless explicitly set in due string e.g.
	// '48 hours'     is future (implicit)
	// 'in 48 hours'  is future (explicit)
	// '48 hours ago' is past (explicit)
	if t, err := naturaldate.Parse(due, currentTime, naturaldate.WithDirection(naturaldate.Future)); err != nil {
		return Reminder{}, err
	} else {
		return Reminder{
			DiscordId: discordId,
			Due:       t,
			Text:      text,
		}, nil
	}
}

func PrintFriendlyTime(t time.Time) string {
	return t.Format("Mon, 02 Jan 2006 15:04:05 MST")
}
