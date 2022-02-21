# Luna

Luna is a discord bot.

### Features
- Announces the person who currently has Pokerus,every 15 minutes, in a channel called #rus-alert. To use Luna in your server, make sure your server has a channel #rus-alert, or she won't annouce anything!
- Tell Luna your PFQ and she will ping you when you have Pokerus: `@luna @example#1234 is DrWho` (the @ is a Discord ping)
- Create polls: `@luna poll: <question>? answer1 / answer2`. Luna replies sends a message where you can vote with emojis. You can also leave out the answers, for example `@luna poll: <question>?` and Luna will default to the answers yes / no.

### Development
Development takes place on the development branch. In development mode, Luna does not make real HTTP requests. You must add two environment variables: 
- `DISCORD_TOKEN` with your discord bot token
- `LUNA_ENV` with the value `development` to turn on development mode. See `examples/launch.json` for a Visual Studio Code launch.json file.

### Production
Note that you must get permission from Pokefarm to use their API. In production, you must set the environment variable `DISCORD_TOKEN`. The main script is `bot.py`.