package stuff

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

const pfq = "https://pokefarm.com"

type Username string

type User struct {
	Name string // name with case and spaces preserved
	Url  string
}

// Create a User from a relative URL
func NewUser(relativeUrl string) (User, error) {
	i := strings.LastIndex(relativeUrl, "/")
	if i < 0 {
		return User{}, fmt.Errorf("%s is not a relative URL", relativeUrl)
	}
	name := relativeUrl[i+1:]
	url := pfq + relativeUrl
	return User{name, url}, nil
}

func NewUsername(name string) Username {
	s := name
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "+")
	return Username(s)
}

func UserExists(name Username) (bool, error) {
	resp, err := client.Head(pfq + "/user/" + string(name))
	if err != nil {
		return false, err
	}
	return resp.StatusCode < 300, nil
}

// Returns the relative URL of the current Pokerus holder e.g. /user/SYSTEM
func Pokerus() (User, error) {
	resp, err := client.Head(pfq + "/user/~pkrs")
	if err != nil {
		return User{}, err
	}
	location, ok := resp.Header["Location"]
	if !ok {
		return User{}, fmt.Errorf("location header not found in response")
	}
	user, err := NewUser(location[0])
	if err != nil {
		return User{}, err
	}
	return user, nil
}



// Caveat: If you restart Luna during the middle of a Pokerus minute,
// this channel will send the same Pokerus holder twice
func PokeursChannel(ch chan User) {
	lastMinute := -1
	for {
		now := time.Now()
		// run every 15 minutes
		if now.Minute()%15 == 0 && now.Minute() != lastMinute && now.Second() > 3 {
			pokerus, err := Pokerus()
			if err != nil {
				// something went wrong, check again in 15 seconds
				log.Printf("error getting pokerus holder %v\n", err)
			} else {
				// send the pokeurs holder to the chhanel
				log.Printf("pokerus holder is %v", pokerus)
				ch <- pokerus
				lastMinute = now.Minute()
			}
		}
		time.Sleep(15 * time.Second)
	}
}
