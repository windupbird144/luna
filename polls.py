"""
poll: where do you live? north america / south america / africa / europe / asia / australia / antarctica 
"""
import re

poll_pattern = re.compile("^poll:\s+(.+?)\?(.*)$")

def parse(text):
    m = poll_pattern.match(text)
    if m is None:
        return
    question = m.group(1) + "?"
    tmp = m.group(2)
    answers = [s.strip() for s in tmp.split("/")] if tmp else ["yes","no"]
    return {
        "question": question,
        "answers": answers
    }