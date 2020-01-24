import json
import logging
import os
from datetime import datetime

import paho.mqtt.client as mqtt


class Sciler:
    """
    Class SccLib sets up the connection and the right handler
    """

    def __init__(self, config, device):
        """
        Initialize device with its configuration json file and python script.
        :param config: configuration json file
        :param device: self of device
        """
        self.device = device
        self.config = json.load(config)
        self.name = self.config.get("id", "not-existent")
        self.host = self.config.get("host", "not-existent")
        self.port = self.config.get("port", "not-existent")
        self.labels = self.config.get("labels", "not-existent")
        if [self.host, self.name, self.port, self.labels].count("not-existent") != 0:
            raise Exception(
                "your config file is missing attributes, "
                "make sure it contains all keys from the manual"
            )

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
        :param level: level of message
        :param buf: message to be logged
        """
        msg = self.name, ", broker log: ", level, ", ", buf
        logging.info(msg)

    def log(self, level, msg):
        """
        log method to make a log with the logging packege
        :param level: level of message
        :param msg: message to be logged
        """
        logging.log(level=level, msg=msg)

    def start(self, loop=None, stop=None):
        """
        Starting method to call from the starting script.
        :param loop: possible event loop
        :param stop: possible end function
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
        A connection status false is send to the back-end
        mqtt client is disconnected and logging is turned off
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
        """
        send_message publishes messages to mqtt topics
        :param topic: topic to publish to
        :param json_message: message to publish in json formatting
        """
        #
        self.client.publish(topic, json_message, 1)
        message_type = topic + " message published"
        logging.info((message_type, json_message))

    def __connect(self):
        """
        Connect method to set up the connection to the broker.
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
        :param client: mqtt client
        :param userdata: the private user data as set in Client() or userdata_set()
        :param flags: response flags sent by the broker
        :param rc:  the connection result
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
        param client: mqtt client
        :param userdata: the private user data as set in Client() or userdata_set()
        :param rc: the connection result
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
        :param msg: message with all component status
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
        :param client: mqtt client
        :param userdata: the private user data as set in Client() or userdata_set()
        :param message: message from the back-end
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
        Correct mapper is called with the content of the message.
        Send confirmation message
        :param message: message from the back-end
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
        """
        check_Message executes all instructions in a message
        :param contents: contents list of instruction from the original instruction message
        """
        for action in contents:
            instruction = action.get("instruction")
            if instruction == "test":
                self.__do_test()
            elif instruction == "status update":
                self.__do_status_update()
            elif instruction == "reset":
                self.__do_reset()
            else:
                success = self.__do_custom_instruction(action)
                if not success:
                    return False
        return True

    def __do_test(self):
        """
        Method that performs instruction `test`
        """
        self.device.test()
        logging.info("instruction performed test")

    def __do_status_update(self):
        """
        Method that performs instruction `status update`
        """
        msg_dict = {
            "device_id": self.name,
            "time_sent": datetime.now().strftime("%d-%m-%Y %H:%M:%S"),
            "type": "connection",
            "contents": {"connection": True},
        }
        msg = json.dumps(msg_dict)
        self.__send_message("back-end", msg)
        self.status_changed()
        logging.info("instruction performed status update")

    def __do_reset(self):
        """
        Method that performs instruction `reset`
        """
        self.device.reset()
        logging.info("instruction performed reset")

    def __do_custom_instruction(self, action):
        """
        Methods that performs custom instruction
        :param action: dictionary with action details such as instruction, component_id and value
        :return: boolean whether the instruction was performed successfully
        """
        success = self.device.perform_instruction(action)
        if success:
            logging.info(("instruction performed", action.instruction))
        else:
            logging.warning(
                (
                    "instruction: " + action.instruction + " could not be performed",
                    action,
                )
            )
            return False

    def __subscribe_topic(self, topic):
        """
        Method to call to subscribe to a topic which the
        sciler system wants to receive from the broker.
        :param topic: topic to subscribe too
        """
        self.client.subscribe(topic=topic)
        logging.info(("subscribed to topic", topic))
