CREATE TABLE usernames (
    discord varchar(200),
    username varchar(200),
    UNIQUE (discord,username)
);

CREATE TABLE reminders (
    guild_id VARCHAR(200),
    discord_id VARCHAR(200),
    due TIMESTAMP,
    text TEXT
)