"""
Environment variables:
    DISCORD_TOKEN (string, required)
"""
import asyncio
import os
import re
from asyncio import events
from datetime import datetime, timedelta
from time import sleep, time

import discord
from discord.ext import commands, tasks

import mappings
import pokerus

pattern = re.compile("^@luna @\S+ is (.+)$")

class Luna(discord.Client):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.announce_pokerus.start()

    async def on_ready(self):
        print('Luna is connected to Discord!')

    @tasks.loop(minutes=15)
    async def announce_pokerus(self):
        name,url = pokerus.get_rus_holder()
        member = mappings.get_member_id(pokerus.normalize(name))
        for guild in self.guilds:
            holder = f"<@!{member}>" if member else name
            reply = f"{holder} has Pokérus <{url}>"
            channel = discord.utils.get(guild.text_channels, name='rus-alert')
            if channel is not None:
                await channel.send(reply)


    @announce_pokerus.before_loop
    async def before_my_task(self):
        await asyncio.sleep(pokerus.seconds_until_next_check())

    async def on_message(self, message):
        # we don't want the bot to reply to itself
        if self.user.id == message.author.id:
            return
        
        # somebody mentioned luna
        if len(message.raw_mentions) and message.raw_mentions[0] == self.user.id:
            m = pattern.match(message.clean_content)
            if len(message.raw_mentions) == 2 and m is not None:
                [_, member_id] = message.raw_mentions
                [_,member] = message.mentions
                pfq_name = pokerus.normalize(m.group(1))
                mappings.add_mapping(pfq_name, member_id)
                await message.channel.send(f"thank you, i added {member.display_name} as {pfq_name}")

if __name__ == "__main__":
    token = os.getenv('DISCORD_TOKEN')
    client = Luna()
    client.run(token)