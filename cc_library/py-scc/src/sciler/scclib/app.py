from datetime import datetime
import os
import json
import paho.mqtt.client as mqtt
import logging


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
        self.labels = self.config.get("labels")

        if not os.path.exists("logs"):
            os.mkdir("logs")
        filename = "logs/log-" + datetime.now().strftime("%d-%m-%YT--%H-%M-%S") + ".txt"
        logging.basicConfig(
            level=logging.INFO,
            format="%(asctime)s [%(levelname)-5.5s]  %(message)s",
            handlers=[logging.FileHandler(filename=filename), logging.StreamHandler()],
        )
        msg = "Start of log for device: " + self.name
        logging.info(msg=msg)

        self.client = mqtt.Client(self.name)
        self.client.on_message = self.__on_message
        self.client.on_log = self.__on_log
        self.client.on_connect = self.__on_connect
        self.client.on_disconnect = self.__on_disconnect
        msg_dict = {
            "device_id": self.name,
            "time_sent": datetime.now().strftime("%d-%m-%Y %H:%M:%S"),
            "type": "connection",
            "contents": {"connection": False},
        }
        msg = json.dumps(msg_dict)
        self.client.will_set("back-end", msg, 1)

    def __on_log(self, level, buf):
        """
        MQTT Client method.
        Broker logger that logs everything happening with the mqtt client.
        """

        msg = self.name, ", broker log: ", level, ", ", buf
        logging.info(msg)

    def log(self, level, msg):
        logging.log(level=level, msg=msg)

    def start(self, loop=None, stop=None):
        """
        Starting method to call from the starting script.
        """
        self.__connect()
        try:
            if loop:
                self.client.loop_start()
                loop()
            else:
                self.client.loop_forever()
        except KeyboardInterrupt:
            logging.info("program was terminated from keyboard input")
        finally:
            if stop:
                stop()
            if loop:
                self.client.loop_stop()
            self.__stop()

    def __stop(self):
        """
        Stop method to call from the starting script.
        """
        msg_dict = {
            "device_id": self.name,
            "time_sent": datetime.now().strftime("%d-%m-%Y %H:%M:%S"),
            "type": "connection",
            "contents": {"connection": False},
        }
        msg = json.dumps(msg_dict)
        self.__send_message("back-end", msg)
        logging.info("cleanly exited ControlBoard program and client")
        self.client.disconnect()
        logging.shutdown()

    def __send_message(self, topic, json_message):
        # TODO what to do when publish fails
        self.client.publish(topic, json_message, 1)
        message_type = topic + " message published"
        logging.info((message_type, json_message))

    def __connect(self):
        """
        Connect method to set up the connection to the broker.
        When connected:
        sends message to topic connection to say its connected,
        subscribes to topic "test"
        starts loop_forever
        """
        try:
            self.client.connect(self.host, self.port, keepalive=10)
            logging.info("connected to broker")
        except ConnectionRefusedError:
            logging.error("connection was refused")
        except TimeoutError:
            logging.error("connecting failed, socket timed out")

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
            for label in self.labels:
                self.__subscribe_topic(label)
            self.__subscribe_topic("client-computers")
            self.__subscribe_topic(self.name)
            msg_dict = {
                "device_id": self.name,
                "time_sent": datetime.now().strftime("%d-%m-%Y %H:%M:%S"),
                "type": "connection",
                "contents": {"connection": True},
            }
            msg = json.dumps(msg_dict)
            self.__send_message("back-end", msg)
            self.status_changed()
            logging.info("connected OK")
        else:
            logging.error(("bad connection, returned code=", rc))
            client.bad_connection_flag = True

    def __on_disconnect(self, client, userdata, rc):
        """
        MQTT Client method.
        When disconnecting from the broker, on_disconnect prints the reason.
        """
        msg_dict = {
            "device_id": self.name,
            "time_sent": datetime.now().strftime("%d-%m-%Y %H:%M:%S"),
            "type": "connection",
            "contents": {"connection": False},
        }
        msg = json.dumps(msg_dict)
        self.__send_message("back-end", msg)
        if rc == 0:
            logging.info(("disconnecting, reason  " + str(rc)))
        else:
            logging.warning("disconnecting, reason " + str(rc))
        client.connected_flag = False
        client.disconnect_flag = True

    def status_changed(self):
        """
        This is called from the client computer to message a status update.
        """
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
            "contents": msg,
        }
        res_msg = json.dumps(json_msg)
        self.__send_message("back-end", res_msg)

    def __on_message(self, client, userdata, message):
        """
        MQTT Client method.
        This method is called when the client receives
         a message from the broken for a subscribed topic.
        The message is printed and send through to the handler.
        """
        logging.info(
            (
                "message received: topic",
                message.topic,
                "message",
                str(message.payload.decode("utf-8")),
            )
        )
        self.__handle(message)

    def __handle(self, message):
        """
        Interpreter of incoming messages.
        Correct sciler mapper is called with the content of the message.
        Send confirmation message
        """
        message = message.payload.decode("utf-8")
        message = json.loads(message)
        if message.get("type") != "instruction":
            logging.warning(
                ("received non-instruction message of type: ", message.get("type"))
            )
        else:
            success = self.__check_message(message.get("contents"))
            conf_msg_dict = {
                "device_id": self.name,
                "time_sent": datetime.now().strftime("%d-%m-%Y %H:%M:%S"),
                "type": "confirmation",
                "contents": {"completed": success, "instructed": message},
            }
            msg = json.dumps(conf_msg_dict)
            self.__send_message("back-end", msg)

    def __check_message(self, contents):
        for action in contents:
            instruction = action.get("instruction")
            if instruction == "test":
                self.device.test()
                logging.info(("instruction performed", instruction))
            elif instruction == "status update":
                msg_dict = {
                    "device_id": self.name,
                    "time_sent": datetime.now().strftime("%d-%m-%Y %H:%M:%S"),
                    "type": "connection",
                    "contents": {"connection": True},
                }
                msg = json.dumps(msg_dict)
                self.__send_message("back-end", msg)
                self.status_changed()
                logging.info(("instruction performed", instruction))
            elif instruction == "reset":
                self.device.reset()
                logging.info(("instruction performed", instruction))
            else:
                (success, failed_action) = self.device.perform_instruction(
                    action
                )  # TODO: remove failed_action as return argument
                if success:
                    logging.info(("instruction performed", instruction))
                else:
                    logging.warning(
                        (
                            "instruction: " + failed_action + " could not be performed",
                            action,
                        )
                    )
                    return False
        return True

    def __subscribe_topic(self, topic):
        """
        Method to call to subscribe to a topic which the
        sciler system wants to receive from the broker.
        """
        self.client.subscribe(topic=topic)
        logging.info(("subscribed to topic", topic))
