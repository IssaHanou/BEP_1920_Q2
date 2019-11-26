import paho.mqtt.client as mqtt
import json
from datetime import datetime


def on_connect(self, userdata, flags, rc):
    if rc == 0:
        self.client.connected_flag = True  # set flag
        print("connected OK")
    else:
        print("Bad connection Returned code=", rc)
        self.client.bad_connection_flag = True


def on_disconnect(client, userdata, rc):
    print("disconnecting reason  " + str(rc))
    client.connected_flag = False
    client.disconnect_flag = True


def on_log(client, userdata, level, buf):
    print("log: ", buf)


class Connection:
    def __init__(self, config, app):
        self.config = config
        self.app = app

        self.name = self.config.get("id")
        self.info = self.config.get("info")
        self.host = self.config.get("host")
        self.client = mqtt.Client(self.name)
        self.client.on_message = self.on_message
        self.client.on_log = on_log
        self.client.on_connect = on_connect
        self.client.on_disconnect = on_disconnect

        print(self)

    def on_message(self, client, userdata, message):
        print("message received ", str(message.payload.decode("utf-8")))
        print("message topic=", message.topic)
        self.app.handler.handle(self, message)

    def connect(self):
        while True:
            try:
                self.client.connect(self.host, 1883, keepalive=60)
                print("we zijn connected")
                msg_dict = {
                    "messageConnectionConfirmation": {
                        "device_id": self.name,
                        "time_sent": datetime.now().strftime("%Y-%m-%dT%H:%M:%S.%f%z"),
                        "type": "connection",
                        "message": {"connection": True},
                    }
                }
                msg = json.dumps(msg_dict)
                self.client.publish("connection", msg)
                self.app.subscribe_topic("test")
                self.client.loop_forever()
                break
            except (ConnectionRefusedError, ConnectionError):
                print("alles is kapot")
