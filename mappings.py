"""
Manages mappings between Pokefarm Q usernames and discord
member IDs.
"""
from datetime import datetime
import sqlite3

con = sqlite3.connect('luna.db')

cur = con.cursor()


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


def add_reminder(person, task, due):
    cur = con.cursor()
    cur.execute('''INSERT INTO reminders (person, task, due) VALUES (?, ?, ?)''', [person, task, due])
    con.commit()

def get_due_reminders(reference_time) -> list:
    """Returns all reminders that are due relative to reference_time"""
    cur = con.cursor()
    cur.execute('''SELECT person, task, due from reminders where due < ?''', [reference_time])
    due_reminders = cur.fetchall()
    cur.close()
    return due_reminders

def delete_due_reminders(reference_time):
    """Delets all reminders that are due relative to reference_time"""
    cur = con.cursor()
    cur.execute('''DELETE from reminders where due < ?''', [reference_time])
    cur.close()
    con.commit()
