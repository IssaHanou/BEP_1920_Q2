from datetime import datetime
import json

import paho.mqtt.client as mqtt


def on_connect(self, userdata, flags, rc):
    if rc == 0:
        self.connected_flag = True  # set flag
        print("connected OK")
    else:
        print("Bad connection Returned code=", rc)
        self.bad_connection_flag = True


def on_disconnect(client, userdata, rc):
    print("disconnecting reason  " + str(rc))
    client.connected_flag = False
    client.disconnect_flag = True


def on_log(client, userdata, level, buf):
    print("log: ", buf)


class App:
    """
    Class App sets up the connection and the right handler
    """

    def __init__(self, config, device):
        self.device = device
        self.config = json.load(config)
        self.name = self.config.get("id")
        self.info = self.config.get("info")
        self.host = self.config.get("host")
        self.client = mqtt.Client(self.name)
        self.client.on_message = self.on_message
        self.client.on_log = on_log
        self.client.on_connect = on_connect
        self.client.on_disconnect = on_disconnect

    def start(self):
        self.connect()

    def subscribe_topic(self, topic):
        print("subscribed to topic", topic)
        self.client.subscribe(topic=topic)

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
        self.client.publish("status", str(jsonmsg))

    def on_message(self, client, userdata, message):
        print("message received ", str(message.payload.decode("utf-8")))
        print("message topic=", message.topic)
        self.handle(self, message)

    def connect(self):
        while True:
            try:
                self.client.connect(self.host, 1883, keepalive=60)
                print("we zijn connected")
                msg_dict = {
                    "messageConnectionConfirmation": {
                        "device_id": self.name,
                        "time_sent":
                            datetime.now().strftime("%Y-%m-%dT%H:%M:%S.%f%z"),
                        "type": "connection",
                        "message": {"connection": True},
                    }
                }
                msg = json.dumps(msg_dict)
                self.client.publish("connection", msg)
                self.subscribe_topic("test")
                self.client.loop_forever()
                break
            except (ConnectionError):
                print("alles is kapot")

    def handle(self, message):
        message = message.payload.decode("utf-8")
        message = json.loads(message)
        message_type = message.get("messageInstructionTest").get("type")
        if message_type == "instruction":
            self.device.incoming_instruction(
                message.get("messageInstructionTest").get("contents")
            )
        elif message_type == "status":
            status = self.device.incoming_status()
            self.send_status_message(status)
