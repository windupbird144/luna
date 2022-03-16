package stuff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("/user/system")
	assert.Equal(t, err, nil)
	assert.Equal(t, user.Name, "system")
	assert.Equal(t, user.Url, "https://pokefarm.com/user/system")
	_, err = NewUser("")
	assert.Error(t, err)
}

func TestNormalize(t *testing.T) {
	assert.Equal(t, string(NewUsername(" hello ")), "hello")
	assert.Equal(t, string(NewUsername("HELLO")), "hello")
	assert.Equal(t, string(NewUsername("h e llo")), "h+e+llo")
}

func TestOnce(t *testing.T) {
	UserExists("system")
}
