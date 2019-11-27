from datetime import datetime
import json

from src.device.base_handler import BaseHandler
from src.device.connection import Connection


class App:
    """
    Class App sets up the connection and the right handler
    """

    connection_class = Connection
    fallback_handler_class = BaseHandler

    connection = None

    def __init__(self, config, device):
        # self.handler_dict = handler_dict or handler_registry.copy()
        self.device = device
        self.config = json.load(config)
        self.name = self.config.get("id")

        self.connection = self.connection_class(config=self.config, app=self)
        # self.verbosity = int(self.config.get('verbosity', 1))
        # self.api_url = config.get('api_url')
        self.handler = BaseHandler

    def start(self):
        self.connection.connect()

    def subscribe_topic(self, topic):
        print("subscribed to topic", topic)
        self.connection.client.subscribe(topic=topic)

    def send_status_message(self, msg):
        jsonmsg = {
            "messageStatusComponent": {
                "device_id": self.name,
                "time_sent": datetime.now().strftime("%Y-%m-%dT%H:%M:%S.%f%z"),
                "type": "status",
                "contents": eval(msg),
            }
        }
        print(jsonmsg)
        self.connection.client.publish("status", str(jsonmsg))
