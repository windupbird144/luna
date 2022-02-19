"""
Manages mappings between Pokefarm Q usernames and discord
member IDs.
"""

import sqlite3

con = sqlite3.connect('luna.db')


def add_mapping(pfq_name, member_id):
    """Adds the mapping pokefarm_name -> discord_name"""
    cur = con.cursor()
    cur.execute('''INSERT INTO mappings (pfq_name, member_id)
    VALUES (?, ?)''', (pfq_name, member_id))
    cur.close()
    con.commit()

def get_member_id(pfq_name):
    """Returns the discord user ID for the pokefarm user if such
    a mapping exists. otherwise"""
    cur = con.cursor()
    cur.execute('''SELECT member_id from mappings
        WHERE pfq_name = ?''', (pfq_name,))
    next = cur.fetchone()
    cur.close()
    if next:
        return next[0]
    return None