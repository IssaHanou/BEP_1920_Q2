import json

from src.device.base_handler import BaseHandler
from src.device.connection import Connection


class App:
    connection_class = Connection
    fallback_handler_class = BaseHandler

    connection = None

    def __init__(self, config, device):
        # self.handler_dict = handler_dict or handler_registry.copy()
        self.device = device
        self.config = json.load(config)
        self.connection = self.connection_class(config=self.config, app=self)
        # self.verbosity = int(self.config.get('verbosity', 1))
        # self.api_url = config.get('api_url')
        self.handler = BaseHandler

    def on_message(self, userdata, message):
        print("message received ", str(message.payload.decode("utf-8")))
        print("message topic=", message.topic)
        self.handler.handle(message)

    def start(self):
        self.connection.connect()

    def subscribe_topic(self, topic):
        print("subscribed to topic", topic)
        self.connection.client.subscribe(topic="test")
        msg_dict = {
            "messageInstructionTest": {
                "device_id": "controlBoard",
                "time_sent": "01-10-2019 13:59:02",
                "type": "instruction",
                "contents": {"instruction": "test"},
            }
        }
        msg = json.dumps(msg_dict)

        self.connection.client.publish("test", msg)
