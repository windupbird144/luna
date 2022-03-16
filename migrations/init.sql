CREATE TABLE usernames (
    discord varchar(200),
    username varchar(200),
    UNIQUE (discord,username)
);