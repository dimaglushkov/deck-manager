import json
import random

import requests
import os

env_file = ".env"
env_vars = dict()
valid_card_codes = [""] * 52
url = "http://127.0.0.1:{}/{}"


def init_env_var():
    global env_vars
    if not os.path.exists(env_file):
        print(f"No {env_file} file")
        exit(1)
    with open(env_file, "r") as ef:
        for line in ef.readlines():
            name, val = line.split("=")
            env_vars[name] = val


def init_valid_cards():
    global valid_card_codes
    ind = 0
    for s in "CDHS":
        for v in range(2, 11):
            valid_card_codes[ind] = str(v) + s
            ind += 1
        for v in "JQKA":
            valid_card_codes[ind] = v + s
            ind += 1


def test_create_default():
    global valid_card_codes
    resp = requests.get(url.format(env_vars["APP_PORT"], "create"))
    deck_data = json.loads(resp.text)

    cards = json.loads(
        requests.get(url.format(env_vars["APP_PORT"], "open") + f"?deck_id={deck_data['deck_id']}").text)
    assert valid_card_codes == [card["code"] for card in cards["cards"]]


def test_create_shuffle():
    deck1_data = json.loads(requests.get(url.format(env_vars["APP_PORT"], "create") + "?shuffled=true").text)
    deck2_data = json.loads(requests.get(url.format(env_vars["APP_PORT"], "create") + "?shuffled=t").text)

    deck1_cards = json.loads(requests.get(url.format(env_vars["APP_PORT"], "open") + f"?deck_id={deck1_data['deck_id']}").text)["cards"]
    deck2_cards = json.loads(requests.get(url.format(env_vars["APP_PORT"], "open") + f"?deck_id={deck2_data['deck_id']}").text)["cards"]

    assert deck1_cards != deck2_cards


def test_create_partial_shuffle():
    partials = "2C,3C,5C,10C,AH"
    deck1_data = json.loads(requests.get(url.format(env_vars["APP_PORT"], "create") + f"?shuffled=t&cards={partials}").text)
    deck2_data = json.loads(requests.get(url.format(env_vars["APP_PORT"], "create") + f"?shuffled=t&cards={partials}").text)

    deck1_cards = json.loads(requests.get(url.format(env_vars["APP_PORT"], "open") + f"?deck_id={deck1_data['deck_id']}").text)["cards"]
    deck2_cards = json.loads(requests.get(url.format(env_vars["APP_PORT"], "open") + f"?deck_id={deck2_data['deck_id']}").text)["cards"]

    assert deck1_cards != deck2_cards
    for card in deck1_cards:
        assert card["code"] in partials
    for card in deck2_cards:
        assert card["code"] in partials


def test_draw_cards():
    global valid_card_codes
    deck_data = json.loads(requests.get(url.format(env_vars["APP_PORT"], "create")).text)
    drawn = 0
    remaining = deck_data["remaining"]

    while remaining > 0:
        to_draw = random.randint(1, remaining)
        cards = json.loads(requests.get(url.format(env_vars["APP_PORT"], "draw") + f"?deck_id={deck_data['deck_id']}&num={to_draw}").text)
        card_codes = [card["code"] for card in cards]
        valid = valid_card_codes[drawn:drawn+to_draw]
        assert valid == card_codes
        drawn += to_draw
        remaining -= to_draw

    deck_data = json.loads(requests.get(url.format(env_vars["APP_PORT"], "open") + f"?deck_id={deck_data['deck_id']}").text)
    assert deck_data["remaining"] == 0


def main():
    init_valid_cards()
    init_env_var()

    test_create_default()
    test_create_shuffle()
    test_create_partial_shuffle()
    test_draw_cards()
    print("OK")


if __name__ == "__main__":
    main()
