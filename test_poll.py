import polls
import pytest

def test_parse_poll_simple():
    my_poll = polls.parse("poll: favorite color? black / pink / red")
    assert my_poll["question"] == "favorite color?"
    assert my_poll["answers"] == ["black", "pink", "red"]