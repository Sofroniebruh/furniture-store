import json
import os

from dotenv import load_dotenv
from flask import Flask, jsonify
import pika
import requests

app = Flask(__name__)
load_dotenv()

RABBITMQ_USER = os.getenv("RABBITMQ_USER")
RABBITMQ_PASSWORD = os.getenv("RABBITMQ_PASSWORD")
API_GATEWAY_URL = os.getenv("AWS_API_GATEWAY")
API_GATEWAY_KEY = os.getenv("AWS_API_GATEWAY_KEY")
RABBIT_PORT = os.getenv("RABBIT_PORT")


def main():
    credentials = pika.PlainCredentials(RABBITMQ_USER, RABBITMQ_PASSWORD)
    conn = pika.BlockingConnection(pika.ConnectionParameters('rabbitmq', RABBIT_PORT, '/', credentials))
    channel = conn.channel()
    channel.queue_declare(queue='account_queue')

    def callback(ch, method, properties, body):
        decoded_body = body.decode('utf-8')
        reply_to = properties.reply_to
        correlation_id = properties.correlation_id

        try:
            data = send_data(decoded_body)

            print("Data: ", data)

            status_code = data.get('statusCode')
            data_body = data.get('body')

            response_body = json.dumps({
                "status_code": status_code,
                "data": data_body,
            })
        except Exception as e:
            response_body = json.dumps({
                "error": str(e)
            })

        ch.basic_publish(
            exchange='',
            routing_key=reply_to,
            properties=pika.BasicProperties(
                correlation_id=correlation_id,
                content_type='application/json'
            ),
            body=response_body.encode('utf-8')
        )

    channel.basic_consume(queue='account_queue', auto_ack=True, on_message_callback=callback)
    channel.start_consuming()


def send_data(decoded_body):
    headers = {
        'Content-Type': 'application/json',
        "x-api-key": API_GATEWAY_KEY,
    }
    url = API_GATEWAY_URL + "/process-data-to-ses"

    response = requests.post(url, headers=headers, json=decoded_body)

    return response.json()


if __name__ == '__main__':
    main()
