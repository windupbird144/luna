# Luna the friendly Discord bot

Luna is a friendly Discord bots with features for the online game Pokefarm Q.

Commands:
- `/hug @user` hug a user
- `/pokerus mypfqusername` tell luna your pfq username and get pings when you have pokerus. If your username changes, run this command again with your new username.
- `/hyperbeam @user` blast a hyperbeam at a user

Automatic features:
- If your guild has a channel #rus-alert, Luna announces the Pokérus host every 15 minutes

## Testing luna

Tests are run against a real database. Run the following commands

1. Create a docker container bound to port 24019

`docker container run -p 24019:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -d --name luna-db postgres`

2. Apply the migrations


`migrate -source "file://./migrations" -database "postgres://postgres@localhost:24019/postgres?sslmode=disable" up`

3. Run the tests

`go test luna/operations`

## Running luna

- Create a postgres database and apply migrations/init.sql
- Build the project with `go build`
- Run with the appropriate [command line flags](./main.go)

## Developing luna
- To modify an existing command locate the files `./operations/$command.go` and `./operations/$command_test.go` and modify them accordingly.
- To create a new command create the files `./operations/$command` and `./operations/$command_test.go` where `$command` should match the Discord command that is executed. Register the new command in `./main.go`
- To create a migration, use (example): `~/.go/bin/migrate create -dir migrations -ext sql -seq 6 add_settings`
- To run Luna `go run main.go -db "postgres://postgres@localhost:24019/postgres?sslmode=disable" -app MY_APP -token MY_TOKEN -pokerusserver ":12345"`