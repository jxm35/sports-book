import asyncio
import datetime
import json
import os

from typing import List, Dict, Any

import aiohttp
from understat import Understat
import boto3


def write_to_queue(results_to_send: List[Dict[str, Any]]) -> None:
    # Create SQS client
    sqs = boto3.client('sqs')

    queue_url = os.environ['SQS_QUEUE_URL']
    response = sqs.send_message(
        QueueUrl=queue_url,
        DelaySeconds=10,
        MessageAttributes={
            'Title': {
                'DataType': 'String',
                'StringValue': 'new_results'
            },
            'Author': {
                'DataType': 'String',
                'StringValue': 'result_scraper'
            }
        },
        MessageBody=(
            json.dumps(results_to_send)
        )
    )
    print(response)


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


def handle_results(event, context):
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)

    results = loop.run_until_complete(get_results())
    filtered_results = [result for result in results if result_was_today(result)]
    if len(filtered_results) == 0:
        print("no results found")
        return "no results sent"

    write_to_queue(filtered_results)
    print("sent {count} results to the queue".format(count=len(filtered_results)))

    return "{count} results sent".format(count=len(filtered_results))
