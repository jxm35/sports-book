import asyncio
import datetime
import json
import logging
import os

from typing import List, Dict, Any

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
                'StringValue': 'New_Results'
            },
            'Author': {
                'DataType': 'String',
                'StringValue': 'result_scraper'
            }
        },
        MessageBody=(
            json.dumps(fixtures_to_send)
        )
    )
    logging.debug(response)


async def get_results() -> List[Dict[str, Any]]:
    async with aiohttp.ClientSession() as session:
        understat = Understat(session)
        results = await understat.get_league_results(
            "epl",
            2023,
        )
        return results


def result_was_today(fixture: Dict[str, Any]) -> bool:
    date = fixture["datetime"]
    dt = datetime.datetime.strptime(date, "%Y-%m-%d %H:%M:%S")
    return dt.date() == datetime.datetime.today().date()


if __name__ == "__main__":
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)

    results = loop.run_until_complete(get_results())
    filtered_results = [result for result in results if result_was_today(result)]
    if len(filtered_results) == 0:
        logging.debug("no results found")
        exit(0)

    write_to_queue(filtered_results)
    logging.debug("sent {count} fixtures to the queue".format(count=len(filtered_results)))

    print(filtered_results)
