import polls
import pytest

def test_parse_poll_simple():
    my_poll = polls.parse("poll: favorite color? black / pink / red")
    assert my_poll["question"] == "favorite color?"
    assert my_poll["answers"] == ["black", "pink", "red"]

def test_parse_poll_yes_no():
    my_poll = polls.parse("poll: can you do empty polls?")
    assert my_poll["answers"] == ["yes", "no"]

def test_parse_poll_not_a_poll():
    assert polls.parse("") is None
    assert polls.parse(":?") is None
    assert polls.parse("poll: hello") is None
    assert polls.parse("poll hello?") is None
