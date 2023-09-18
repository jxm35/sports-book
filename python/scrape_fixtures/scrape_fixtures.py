import asyncio
import datetime
import json
import logging
import os

from typing import List, Union, Dict, Any

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
    print(response)


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


def handle_fixtures(event, context):
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
    print("sent {count} fixtures to the queue".format(count=len(filtered_fixtures)))

    return "{count} fixtures sent".format(count=len(filtered_fixtures))


# if __name__ == "__main__":
#     loop = asyncio.new_event_loop()
#     asyncio.set_event_loop(loop)
#
#     fixtures = loop.run_until_complete(get_fixtures())
#     filtered_fixtures = [fixture for fixture in fixtures if fixture_is_today(fixture)]
#     if len(filtered_fixtures) == 0:
#         logging.debug("no fixtures found")
#         exit(0)
#
#     # we have some fixtures today to process
#     # for fixture in filtered_fixtures:
#     #     save_fixture(conn, fixture, year_id)
#     # logging.debug("added {count} fixtures to database".format(count=len(filtered_fixtures)))
#
#     write_to_queue(filtered_fixtures)
#     logging.debug("sent {count} fixtures to the queue".format(count=len(filtered_fixtures)))
#
#     print(filtered_fixtures)
