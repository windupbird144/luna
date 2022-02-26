"""
Environment variables:
    DISCORD_TOKEN (string, required)
"""
import asyncio
import os
import random
import re
from asyncio import events
from datetime import datetime, timedelta, timezone
from time import sleep, time
from venv import create

import discord
from discord.ext import commands, tasks

import env
import mappings
import pokerus
import polls

import dateparser

"""Luna reacts to messages matching one of these regex patterns"""
add_mapping_pattern = re.compile("^@luna @.+? is (.+)$")
create_poll_pattern = re.compile("^@luna (poll: .+)$")
hug_pattern = re.compile("^@luna hug @.+?$")
remindme_pattern = re.compile("^@luna remind me (to .+) (in .+)$")
choose_pattern = re.compile("^@luna choose: (.+)$")

"""The list of emojis used to enumerate the poll answers. The number of
emojis in this list determines the maximum number of distinct answers on a
poll. Entries in this list must be valid Unicode emojis"""
poll_emojis = [
    "\U0001f1e6", # regional indicator A
    "\U0001f1e7", # regional indicator B
    "\U0001f1e8", # ...
    "\U0001f1e9",
    "\U0001f1ea",
    "\U0001f1eb",
    "\U0001f1ec",
    "\U0001f1ed",
    "\U0001f1ee",
    "\U0001f1ef",
    "\U0001f1f0",
    "\U0001f1f1",
    "\U0001f1f2",
    "\U0001f1f3",
    "\U0001f1f4",
    "\U0001f1f5",
    "\U0001f1f6",
    "\U0001f1f7",
    "\U0001f1f8",
    "\U0001f1f9",
    "\U0001f1fa",
    "\U0001f1fb",
    "\U0001f1fc",
    "\U0001f1fd", # ...
    "\U0001f1fe", # regional indicator Y
    "\U0001f1ff"  # regional indicator Z
]

def split_and_strip(s: str, sep: str) -> list[str]:
    return [x.strip() for x in s.split(sep)]

class Luna(discord.Client):
    def parse(self, message):
        """Parses a Discord message to a command. returns None if the message
        is not a command for luna."""
        # do not reply to yourself
        if self.user.id == message.author.id:
            return None
        
        # luna was not pinged
        if not self.user.id in message.raw_mentions:
            return None

        # matching against the 'add mapping' command
        m = add_mapping_pattern.match(message.clean_content)
        if m is not None:
            [_, member_id] = message.raw_mentions
            [_,member] = message.mentions
            pfq_display_name = m.group(1)
            pfq_name = pokerus.normalize(pfq_display_name)
            return 'add_mapping', member_id, pfq_name, member.display_name, pfq_display_name
        
        # match against the poll pattern
        m = create_poll_pattern.match(message.clean_content)
        if m is not None:
            return 'create_poll', m.group(1) 

        # match against the hug pattern
        m = hug_pattern.match(message.clean_content)
        if m is not None and len(message.raw_mentions) == 2:
            return 'hug', message.raw_mentions[1]

        # match against the remindme pattern
        m = remindme_pattern.match(message.clean_content)
        if m is not None:
            person, task, display_time = message.author.id, m.group(1), m.group(2)
            time = dateparser.parse(display_time)
            # must be a naive datetime
            if time is None or time.tzinfo is not None:
                return
            time = time.astimezone(timezone.utc).timestamp()
            return 'remindme', person, task, time, display_time

        # match against choose pattern
        m = choose_pattern.match(message.clean_content)
        if m is not None:
            choices = m.group(1)
            choices = split_and_strip(choices, "/")
            chosen = random.choice(choices)
            return 'choose', chosen
        
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.announce_pokerus.start()
        self.check_reminders.start()

    async def on_ready(self):
        print('Luna is connected to Discord!')

    @tasks.loop(minutes=15)
    async def announce_pokerus(self):
        name,url = pokerus.get_rus_holder()
        member = mappings.get_member_id(pokerus.normalize(name))
        for guild in self.guilds:
            holder = f"<@!{member}>" if member else name
            reply = f"{holder} has Pokérus <{url}>"
            # prefix the reply with a message if it is reset
            now = datetime.utcnow()
            if now.hour == 0 and now.minute == 0:
                reply= f"Happy reset! {reply}"
            channel = discord.utils.get(guild.text_channels, name='rus-alert')
            if channel is not None:
                await channel.send(reply)

    @tasks.loop(seconds=12)
    async def check_reminders(self):
        reference_time = datetime.now().timestamp()
        for guild in self.guilds:
            reminders = mappings.get_due_reminders(reference_time)
            channel = discord.utils.get(guild.text_channels, name='bot')
            for (person, task, due) in reminders:
                msg = f"<@!{person}> here is your reminder {task}"
                await channel.send(msg)
            mappings.delete_due_reminders(reference_time)

    @announce_pokerus.before_loop
    async def before_my_task(self):
        # Return immediately in development
        if env.is_development():
            return await self.wait_until_ready()
        # Wait until the next :00:20, :15:20, :30:20, :45:20
        await asyncio.sleep(pokerus.seconds_until_next_change())

    async def on_message(self, message):
        command = self.parse(message)
        if command is None:
            return
        elif command[0] == 'add_mapping':
            _, member_id, pfq_name, member_display_name, pfq_display_name = command
            user_exists,resolved_url = pokerus.user_exists(pfq_name)
            if user_exists:
                mappings.add_mapping(pfq_name, member_id)
                await message.channel.send(f"thank you, i added {member_display_name} as {pfq_display_name}")
            else:
                apology = f"sorry, i can't find {pfq_display_name} on pfq. i looked at the url <{resolved_url}>"
                has_special_chars = not re.match("^[a-zA-Z0-9 ]+$", pfq_name)
                if has_special_chars:
                    apology = f"{apology}. please try again without special characters?"
                await message.channel.send(apology)
        elif command[0] == 'create_poll':
            _, pollstring = command
            my_poll = polls.parse(pollstring)
            if my_poll is not None:
                q = my_poll["question"]
                answers = my_poll["answers"]
                to_send = f"poll: {q}\n"
                max_answers = len(poll_emojis) # we can't have more answers in the poll than there are react emojis
                answers = answers[:max_answers]
                for i,a in enumerate(answers[:max_answers]):
                    to_send = f"{to_send}{poll_emojis[i]}\t{a}\n"
                reply = await message.channel.send(to_send)
                for i in range(len(answers)):
                    await reply.add_reaction(poll_emojis[i])
        elif command[0] == 'hug':
            _, to_hug = command
            to_send = f"\*hugs <@!{to_hug}>\*"
            await message.channel.send(to_send)
        elif command[0] == 'remindme':
            _, person, task, due, display_time = command
            mappings.add_reminder(person, task, due)
            reply = f"thanks, i will remind you {task} {display_time}"
            await message.reply(reply)
        elif command[0] == 'choose':
            _, choice = command
            to_reply = f"i choose {choice}"
            await message.reply(to_reply)

if __name__ == "__main__":
    token = os.getenv('DISCORD_TOKEN')
    client = Luna()
    client.run(token)
