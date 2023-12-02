package operations

import "database/sql"

func DeleteFromRemindersTableByDiscordId(conn *sql.DB, discordId string) (bool, error) {
	return RemoveUserByDiscordId(conn, discordId)
}
