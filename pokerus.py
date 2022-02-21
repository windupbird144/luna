import re
import env
from datetime import datetime, timedelta

import requests

# Do not use real requests in development
if env.is_development():
    requests = None

base_url = 'https://pokefarm.com'
pkrs_url = '/user/~pkrs'
def to_absolute(path):
    r'''Turns a relative PFQ URL into an absolute URL
    
    Example:
    to_absolute('/farm') == 'https://pokefarm.com/farm'
    '''
    return base_url + path

# returns the relative url of the pokerus holder
def fetch_pokerus_holder_url():
    r'''Make a HTTP request to get the relative URL of
    the current Pokerus holder's profile:

    Example:
    fetch_pokerus_holder_url() == '/user/system'
    '''
    if env.is_development():
        return "/user/system"
    r = requests.head(to_absolute(pkrs_url))
    assert r.status_code == 302
    pokerus_holder_url = r.headers['location']
    return pokerus_holder_url

def parse_user_url(relative_url):
    r'''Parse the relative url of a player's profile
    
    Example:
    name, url = parse_user_url('/user/system')
    name == 'system'
    url == 'https://pokefarm.com/system'
    '''
    name = relative_url.split('/')[-1]
    url = to_absolute(relative_url)
    return name, url


def get_rus_holder():
    return parse_user_url(fetch_pokerus_holder_url())

def user_exists(name):
    if env.is_development():
        return name == "system", "https://pokefarm.com/user/system"
    resolved = to_absolute(f"/user/{name}")
    r = requests.head(resolved)
    return r.status_code < 300, resolved

def normalize(name):
    name = name.strip()
    name = re.sub(r"\s", "+", name)
    name = name.lower()
    return name

def seconds_until_next_change(start : datetime =None):
    """Returns the number of seconds until the next pokerus change. Changes happen
    when the minute is 0, 30, 15 or 45 at the second 20
    
    Parameters:
    start (optional) - A datetime from which to calculate the number of seconds.
    When omitted, the function use datetime.now()
    """
    start = start or datetime.now()
    if start.minute % 15 == 0 and start.second <= 20:
        return 20 - start.second
    else:
        target = start + timedelta(minutes=1)
        while target.minute % 15 != 0:
            target = target + timedelta(minutes=1)
        return (target - start).seconds
