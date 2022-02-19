import os

env = os.getenv("LUNA_ENV")

def is_development():
    return env == "development"