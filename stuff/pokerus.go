package stuff

import (
	"database/sql"
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
