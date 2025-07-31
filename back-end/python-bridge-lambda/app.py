import os

from dotenv import load_dotenv
from flask import Flask
import pika

app = Flask(__name__)
load_dotenv()

RABBITMQ_USER = os.getenv("RABBITMQ_USER")
RABBITMQ_PASSWORD = os.getenv("RABBITMQ_PASSWORD")


def main():
    credentials = pika.PlainCredentials(RABBITMQ_USER, RABBITMQ_PASSWORD)
    conn = pika.BlockingConnection(pika.ConnectionParameters('rabbitmq', 5672, '/', credentials))
    channel = conn.channel()
    channel.queue_declare(queue='verifyEmail')

    def callback(ch, method, properties, body):
        decodedbody = body.decode('utf-8')
        print("Received %r" % decodedbody)

    channel.basic_consume(queue='verifyEmail', auto_ack=True, on_message_callback=callback)
    channel.start_consuming()


if __name__ == '__main__':
    main()
