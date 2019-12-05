from datetime import datetime
import json

import paho.mqtt.client as mqtt

from cc_library.src.sciler.scclib.logger import Logger


class SccLib:
    """
    Class SccLib sets up the connection and the right handler
    """

    def __init__(self, config, device):
        """
        Initialize device with its configuration json file and python script.
        """
        self.device = device
        self.config = json.load(config)
        self.name = self.config.get("id")
        self.info = self.config.get("description")
        self.host = self.config.get("host")
        self.port = self.config.get("port")
        self.logger = Logger()
        self.logger.log("Start of log for device: " + self.name)

        self.statusChanged = self.status_changed
        self.client = mqtt.Client(self.name)
        self.client.on_message = self.__on_message
        self.client.on_log = self.__on_log
        self.client.on_connect = self.__on_connect
        self.client.on_disconnect = self.__on_disconnect

    def __on_log(self, level, buf):
        """
        MQTT Client method.
        Broker logger that logs everything happening with the mqtt client.
        """
        print(self.name, ", broker log: ", buf)

    def start(self):
        """
        Starting method to call from the starting script.
        """
        self.__connect()

    def __send_message(self, topic, json_message):
        # TODO what to do when publish fails
        self.client.publish(topic, json_message)
        message_type = topic + " message published"
        self.logger.log((message_type, json_message))

    def __connect(self):
        """
        Connect method to set up the connection to the broker.
        When connected:
        sends message to topic connection to say its connected,
        subscribes to topic "test"
        starts loop_forever
        """
        while True:
            try:
                self.client.connect(self.host, self.port, keepalive=60)
                self.logger.log("connected to broker")
                msg_dict = {
                    "device_id": self.name,
                    "time_sent": datetime.now().strftime("%d-%m-%Y %H:%M:%S"),
                    "type": "connection",
                    "contents": {"connection": True},
                }

                msg = json.dumps(msg_dict)
                self.__send_message("connection", msg)
                self.__subscribe_topic("client-computers")
                self.__subscribe_topic("test")
                self.__subscribe_topic(self.name)
                self.client.loop_forever()
                break
            except ConnectionRefusedError:
                self.logger.log("ERROR: connection was refused")

    def __on_connect(self, client, userdata, flags, rc):
        """
        MQTT Client method.
        When trying to connect to the broker,
        on_connect will return the result of this action.
        userdata:   the private user data as set in Client() or userdata_set()
        flags:      response flags sent by the broker
        rc:         the connection result
        """
        if rc == 0:
            client.connected_flag = True  # set flag
            self.logger.log("connected OK")
        else:
            self.logger.log(("bad connection, returned code=", rc))
            client.bad_connection_flag = True

    def __on_disconnect(self, client, userdata, rc):
        """
        MQTT Client method.
        When disconnecting from the broker, on_disconnect prints the reason.
        """
        msg_dict = {
            "device_id": id(client),
            "time_sent": datetime.now().strftime("%d-%m-%YT%H:%M:%S"),
            "type": "connection",
            "contents": {"connection": False},
        }

        msg = json.dumps(msg_dict)
        self.__send_message("connection", msg)
        self.logger.log(("disconnecting, reason  " + str(rc)))
        client.connected_flag = False
        client.disconnect_flag = True
        self.logger.close()

    def status_changed(self, channel):
        """
        This is called from the client computer to message a status update.
        """
        self.logger.log("status changed of pin " + str(channel))
        self.__send_status_message(self.device.get_status())

    def __send_status_message(self, msg):
        """
        Method to send status messages to the topic status.
        msg should be a dictionary/json with components
         as keys and its status as value
        """
        json_msg = {
            "device_id": self.name,
            "time_sent": datetime.now().strftime("%d-%m-%Y %H:%M:%S"),
            "type": "status",
            "contents": eval(msg),
        }
        msg = json.dumps(json_msg)
        self.__send_message("status", msg)

    def __on_message(self, client, userdata, message):
        """
        MQTT Client method.
        This method is called when the client receives
         a message from the broken for a subscribed topic.
        The message is printed and send through to the handler.
        """
        self.logger.log(("message received: topic", message.topic, "message", str(message.payload.decode("utf-8"))))
        self.__handle(message)

    def __handle(self, message):
        """
        Interpreter of incoming messages.
        Correct sciler mapper is called with the content of the message.
        Send confirmation message
        """
        message = message.payload.decode("utf-8")
        message = json.loads(message)
        failed_result = self.device.perform_instruction(message.get("contents"))
        msg_dict = {
            "device_id": self.name,
            "time_sent": datetime.now().strftime("%d-%m-%Y %H:%M:%S"),
            "type": "confirmation",
            "contents": {"completed": not failed_result, "instructed": message},
        }
        msg = json.dumps(msg_dict)
        self.__send_message("confirmation", msg)
        if failed_result:
            self.logger.log(("instruction could not be performed", message))
        else:
            self.logger.log(("instruction performed", message))

    def __subscribe_topic(self, topic):
        """
        Method to call to subscribe to a topic which the
        sciler system wants to receive from the broker.
        """
        self.client.subscribe(topic=topic)
        self.logger.log(("subscribed to topic", topic))
