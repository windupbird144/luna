# Luna the friendly Discord bot

Luna is a friendly Discord bots with features for the online game Pokefarm Q.

Commands:
- `/hug @user` hug a user
- `/pokerus mypfqusername` tell luna your pfq username and get pings when you have pokerus

Other features:
- Announces the Pokerus host every 15 minutes in a channel #rus-alert if it exists

To run luna:

- Create a postgres database and apply migrations/init.sql
- Build the project with `go build`
- Run with the appropriate [command line flags](./main.go)
