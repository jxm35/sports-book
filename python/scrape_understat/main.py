import asyncio
import json
import datetime
import logging

import aiohttp

import mysql.connector

from understat import Understat


async def getPlayersInTeam(conn, year):
    async with aiohttp.ClientSession() as session:
        understat = Understat(session)
        players = await understat.get_league_players("serie_a", year)
        for player in players:
            id = int(player['id'])
            name = player['player_name']
            pos = player['position']
            savePlayer(conn=conn, name=name, position=pos, us_id=id)


async def getPlayersInMatch(conn, matchTeamDict, matchesDict, playersDict, home_or_away):
    async with aiohttp.ClientSession() as session:
        understat = Understat(session)
        for match_id, sql_id in matchesDict.items():
            players = await understat.get_match_players(match_id)
            for id, player in players[home_or_away].items():
                player_id = playersDict.get(player['player'])
                if player_id is None:
                    logging.warning("could not find player {player_name}".format(player_name=player['player']))
                    continue
                # player_id = playersDict[player['player']]
                goals = int(player['goals'])
                xG = float(player['xG'])
                minutes = int(player['time'])
                xGChain = float(player['xGChain'])
                xGBuildup = float(player['xGBuildup'])
                assists = int(player['assists'])
                xA = float(player['xA'])
                key_passes = int(player['key_passes'])
                yellow_cards = int(player['yellow_card'])
                red_cards = int(player['red_card'])
                match = sql_id
                team = matchTeamDict[sql_id][player['h_a']]
                saveAppearance(conn=conn, player=player_id, match=match, team=team, goals=goals, expected_goals=xG,
                               expected_goals_chain=xGChain, expected_goals_buildup=xGBuildup, assists=assists,
                               expected_assists=xA, key_passes=key_passes, yellow_cards=yellow_cards,
                               red_cards=red_cards, minutes=minutes, us_id=id)


async def getResults(conn, teamDict, year, yearId):
    async with aiohttp.ClientSession() as session:
        understat = Understat(session)
        results = await understat.get_league_results(
            "serie_a",
            year,
        )
        for result in results:
            id = int(result['id'])
            homeTeam = teamDict[int(result['h']['id'])]
            awayTeam = teamDict[int(result['a']['id'])]
            homeGoals = int(result['goals']['h'])
            awayGoals = int(result['goals']['a'])
            homeXg = float(result['xG']['h'])
            awayXg = float(result['xG']['a'])
            dateString = result['datetime']
            saveMatch(conn=conn, date=dateString, home_team=homeTeam, away_team=awayTeam, home_goals=homeGoals,
                      away_goals=awayGoals, home_expected_goals=homeXg, away_expected_goals=awayXg, us_id=id, yearId=yearId)


async def getTeams(conn, year):
    async with aiohttp.ClientSession() as session:
        understat = Understat(session)
        teams = await understat.get_teams(
            "serie_a",
            year,
        )
        for team in teams:
            id = int(team['id'])
            name = team['title']
            saveTeam(conn=conn, name=name, us_id=id)


"""
competition: id code year
team: id name us_id
player: id name us_id position

// get all teams
// get all players

// iterate through each team, and get their home games
match: id date home_team away_team {competiton} h_g a_g h_xg  a_xg us_id 

// get player stats for game
appearance: id player team {match} g xG xGChain xGBuildup a xA kp yc rc  time us_id
"""
"""
 {
    "h": {"id": "89",
        "title": "Manchester United",
        "short_title": "MUN"}
}
"""


def getConnection():
    mydb = mysql.connector.connect(
        host="127.0.0.1",
        user="root",
        password="password",
        database="sports-book"
    )
    return mydb


def saveTeam(conn, name, us_id):
    mycursor = conn.cursor()
    sql = "INSERT IGNORE INTO `sports-book`.team (name, us_id) VALUES (%s, %s);"
    values = (name, us_id)
    mycursor.execute(sql, values)
    conn.commit()


def savePlayer(conn, name, position, us_id):
    mycursor = conn.cursor()
    sql = "INSERT IGNORE INTO `sports-book`.player (name, position, us_id) VALUES (%s, %s, %s);"
    values = (name, position, us_id)
    mycursor.execute(sql, values)
    conn.commit()


def saveMatch(conn, date, home_team, away_team, home_goals, away_goals, home_expected_goals, away_expected_goals,
              us_id, yearId):
    mycursor = conn.cursor()
    sql = "INSERT INTO `sports-book`.match (date, home_team, away_team, competition, home_goals, away_goals, " \
          "home_expected_goals, away_expected_goals, us_id) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s);"
    values = (date, home_team, away_team, yearId, home_goals, away_goals, home_expected_goals, away_expected_goals, us_id)
    mycursor.execute(sql, values)
    conn.commit()


def saveAppearance(conn, player, match, team, goals, expected_goals, expected_goals_chain, expected_goals_buildup,
                   assists, expected_assists, key_passes, yellow_cards, red_cards, minutes, us_id):
    mycursor = conn.cursor()
    sql = "INSERT INTO `sports-book`.appearance (player, `match`, team, goals, expected_goals, expected_goals_chain, " \
          "expected_goals_buildup, assists, expected_assists, key_passes, yellow_cards, red_cards, minutes, " \
          "us_id) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s);"
    values = (
        player, match, team, goals, expected_goals, expected_goals_chain, expected_goals_buildup, assists,
        expected_assists, key_passes, yellow_cards, red_cards, minutes, us_id)
    mycursor.execute(sql, values)
    conn.commit()


def loadTeamsMap(conn):
    cur = conn.cursor()
    # extract a player list from db
    cur.execute("SELECT id, us_id FROM `sports-book`.team")
    result_set = cur.fetchall()
    teamsDict = {row[1]: row[0] for row in
                 result_set}  # OR {row[0]: row[1] for row in result_set}
    return teamsDict


def loadMatches(conn, yearId):
    cur = conn.cursor()
    # extract a player list from db
    sql = f"SELECT id, us_id FROM `sports-book`.`match` where competition = {yearId}"
    cur.execute(sql)
    result_set = cur.fetchall()
    matchesDict = {row[1]: row[0] for row in
                   result_set}  # OR {row[0]: row[1] for row in result_set}
    return matchesDict


def loadPlayers(conn):
    cur = conn.cursor()
    # extract a player list from db
    cur.execute("SELECT id, name FROM `sports-book`.`player`")
    result_set = cur.fetchall()
    playersDict = {row[1]: row[0] for row in
                   result_set}  # OR {row[0]: row[1] for row in result_set}
    return playersDict


def loadMatchTeams(conn):
    cur = conn.cursor()
    # extract a player list from db
    cur.execute("SELECT id, home_team, away_team FROM `sports-book`.`match`")
    result_set = cur.fetchall()
    playersDict = {row[0]: {"h": row[1], "a": row[2]} for row in
                   result_set}  # OR {row[0]: row[1] for row in result_set}
    return playersDict


if __name__ == "__main__":
    print("started")
    conn = getConnection()
    print("db connection made")
    year = 2022
    yearId = 37

    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)

    loop.run_until_complete(getTeams(conn, year))
    print("teams retrieved")

    loop.run_until_complete(getPlayersInTeam(conn, year))
    print("players retrieved")

    teamsDict = loadTeamsMap(conn)
    loop.run_until_complete(getResults(conn, teamsDict, year, yearId))
    print("results retrieved")

    matchesDict = loadMatches(conn, yearId)
    playersDict = loadPlayers(conn)
    matchTeamsDict = loadMatchTeams(conn)
    print("loaded matches, players, and lineups")
    loop.run_until_complete(getPlayersInMatch(conn, matchTeamsDict, matchesDict, playersDict, 'h'))
    print("home appearances complete")
    loop.run_until_complete(getPlayersInMatch(conn, matchTeamsDict, matchesDict, playersDict, 'a'))
    print("away appearances complete")
    print("finished season")
