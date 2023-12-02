package operations

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteFromRemindersTableByDiscordId(t *testing.T) {
	db, _ := sql.Open("postgres", "host=localhost user=postgres port=24019 sslmode=disable")
	defer db.Close()

	// setup: create a user with id 1
	_, err := db.Exec("insert into usernames (discord, username) values ('1', '')")
	assert.Nil(t, err)

	// delete an id that does not exist
	deleted, err := DeleteFromRemindersTableByDiscordId(db, "-1")
	assert.Nil(t, err)
	assert.False(t, deleted)

	// delete an id that exists
	deleted, err = DeleteFromRemindersTableByDiscordId(db, "1")
	assert.Nil(t, err)
	assert.True(t, deleted)

	// verify it no longer exists
	var i int
	row := db.QueryRow("select 1 from usernames where discord = '1'")
	assert.Equal(t, sql.ErrNoRows, row.Scan(&i))
}
