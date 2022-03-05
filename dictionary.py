import requests

def get_definitions(word : str) -> list[str] or None:
    r = requests.get(f"https://api.dictionaryapi.dev/api/v2/entries/en/{word}")
    if r.status_code >= 300:
        return None
    # the response from the dictionary api is quite complex, simplify it
    # to just definitions. can be changed to use more of the info later
    # if needed.
    definitions = []
    body = r.json()
    for entry in body:
        for meaning in entry['meanings']:
            for definition in meaning['definitions']:
                definitions.append(definition['definition'])
    return definitions