import requests
import csv

# https://api.sofascore.com/api/v1/unique-tournament/17/season/41886/events/round/1
# https://api.sofascore.com/api/v1/unique-tournament/17/season/37036/events/round/2
# https://api.sofascore.com/api/v1/event/10385503/lineups
# https://api.sofascore.com/api/v1/event/10385724/statistics



# Get the data from API
match_id = "10385727"
player_id = "795228"
url = f"https://api.sofascore.com/api/v1/event/{match_id}/player/{player_id}/statistics"
headers = {"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:66.0) Gecko/20100101 Firefox/66.0"}
r = requests.get(url, headers=headers)
json_data = r.json()

player = json_data["player"]
stats = json_data["statistics"]

minutes = stats["minutesPlayed"]
xG = stats["expectedGoals"]
xA = stats["expectedAssists"]
goals = stats["goals"] or 0
assists = stats["goalAssist"] or 0
keyPasses = stats["keyPass"] or 0

print()