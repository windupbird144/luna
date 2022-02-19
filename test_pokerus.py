import pokerus
from datetime import datetime

def test_seconds_until_next_change():
    assert 0 == pokerus.seconds_until_next_change(datetime(2020,1,1,12,15,20,123))
    assert 13 == pokerus.seconds_until_next_change(datetime(2020,1,1,12,30,7))
    assert 15*60 == pokerus.seconds_until_next_change(datetime(2020,1,1,12,30,21))
    assert 7*60 == pokerus.seconds_until_next_change(datetime(2020,1,1,12,38,20))
    assert pokerus.seconds_until_next_change() <= 15*60

def test_normalize():
    assert pokerus.normalize(" hello world ") == "hello+world"
    assert pokerus.normalize("HELLOWORLD") == "helloworld"
