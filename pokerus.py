from datetime import datetime, timedelta
import re
import requests
import schedule

# Example:
#
# curl --head https://pokefarm.com/user/~pkrs
#
# Response:
# Location: /users/lakaihia 

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

def normalize(name):
    name = name.strip()
    name = re.sub(r"\s", "+", name)
    name = name.lower()
    return name

def seconds_until_next_check():
    """Returns the seconds until the next :00:15 :15:15, :30:15, :45:15"""
    for m in ["00","15","30","45"]:
        schedule.every().hour.at(f"{m}:20").do(lambda x: x)
    return schedule.idle_seconds()

if __name__ == '__main__':
    now = datetime.now()
    nxt = seconds_until_next_check()
    print(f"it is {now}")
    print(f"seconds until next check: {nxt}")
    print(f"btw that's {nxt//60} minutes {nxt%60} seconds")