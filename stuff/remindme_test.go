package stuff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tj/go-naturaldate"
)

func TestNewReminder(t *testing.T) {
	// Test if the date parses works as you would expect
	base := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	if tim, err := naturaldate.Parse("in 48 hours", base); err != nil {
		t.Fatal(err)
	} else {
		target := time.Date(2000, time.January, 3, 0, 0, 0, 0, time.UTC)
		assert.Equal(t, tim.Unix(), target.Unix())
	}

}
