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


def main():
    credentials = pika.PlainCredentials(RABBITMQ_USER, RABBITMQ_PASSWORD)
    conn = pika.BlockingConnection(pika.ConnectionParameters('rabbitmq', 5672, '/', credentials))
    channel = conn.channel()
    channel.queue_declare(queue='verifyEmail')

    def callback(ch, method, properties, body):
        decoded_body = body.decode('utf-8')
        send_data(decoded_body)

    channel.basic_consume(queue='verifyEmail', auto_ack=True, on_message_callback=callback)
    channel.start_consuming()


def send_data(decoded_body):
    headers = {
        'Content-Type': 'application/json',
        'x-api-key': API_GATEWAY_KEY,
    }
    payload = json.loads(decoded_body)
    url = API_GATEWAY_URL + "/forward-data-to-email-service"

    response = requests.post(url, headers=headers, json=payload)

    return jsonify({
        "status": response.status_code,
        "body": response.json()
    })


if __name__ == '__main__':
    main()
