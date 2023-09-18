import asyncio
import datetime
import json
import logging
import os

from mysql.connector import MySQLConnection, CMySQLConnection
from mysql.connector.pooling import PooledMySQLConnection
from typing import List, Union, Dict, Any

import mysql.connector
import aiohttp
from understat import Understat
import boto3


def write_to_queue(fixtures_to_send: List[Dict[str, Any]]) -> None:
    # Create SQS client
    sqs = boto3.client('sqs')

    queue_url = os.environ['SQS_QUEUE_URL']
    response = sqs.send_message(
        QueueUrl=queue_url,
        DelaySeconds=10,
        MessageAttributes={
            'Title': {
                'DataType': 'String',
                'StringValue': 'new_fixtures'
            },
            'Author': {
                'DataType': 'String',
                'StringValue': 'fixture_scraper'
            }
        },
        MessageBody=(
            json.dumps(fixtures_to_send)
        )
    )
    logging.debug(response)


async def get_fixtures() -> List[Dict[str, Any]]:
    async with aiohttp.ClientSession() as session:
        understat = Understat(session)
        results = await understat.get_league_fixtures(
            "epl",
            2023,
        )
        return results


def fixture_is_today(fixture: Dict[str, Any]) -> bool:
    date = fixture["datetime"]
    dt = datetime.datetime.strptime(date, "%Y-%m-%d %H:%M:%S")
    return dt.date() == datetime.datetime.today().date()


def getConnection() -> Union[PooledMySQLConnection, MySQLConnection, CMySQLConnection]:
    mydb = mysql.connector.connect(
        host="127.0.0.1",
        user="root",
        password="password",
        database="sports-book"
    )
    return mydb


# def save_fixture(
#         connection: Union[PooledMySQLConnection, MySQLConnection, CMySQLConnection],
#         fixture: Dict[str, Any],
#         season_id: int) -> None:
#     cursor = connection.cursor()
#
#     home_team = fixture['h']['id']
#     away_team = fixture['a']['id']
#     us_id = fixture['id']
#
#     sql = "INSERT INTO `sports-book`.match (date, home_team, away_team, competition, home_goals, away_goals, " \
#           "home_expected_goals, away_expected_goals, us_id) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s);"
#     values = (
#         datetime.datetime.today().date(), home_team, away_team, season_id, -1, -1, -1, -1, us_id)
#     cursor.execute(sql, values)
#     connection.commit()


if __name__ == "__main__":
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)

    fixtures = loop.run_until_complete(get_fixtures())
    filtered_fixtures = [fixture for fixture in fixtures if fixture_is_today(fixture)]
    if len(filtered_fixtures) == 0:
        logging.debug("no fixtures found")
        exit(0)

    # we have some fixtures today to process
    # for fixture in filtered_fixtures:
    #     save_fixture(conn, fixture, year_id)
    # logging.debug("added {count} fixtures to database".format(count=len(filtered_fixtures)))

    write_to_queue(filtered_fixtures)
    logging.debug("sent {count} fixtures to the queue".format(count=len(filtered_fixtures)))
    
    print(filtered_fixtures)
