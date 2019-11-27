import json


# from src.scripts import start_device


class BaseHandler:
    """
    Class for BaseHandler,
    maybe will do the methods of device mapping in a later stage.
    """

    def __init__(self):
        self.status = None

    def on_change(self, status):
        self.status = status

    def handle(self, message):
        message = message.payload.decode("utf-8")
        message = json.loads(message)
        message_type = message.get("messageInstructionTest").get("type")
        if message_type == "instruction":
            self.app.device.incoming_instruction(
                message.get("messageInstructionTest").get("contents")
            )
        elif message_type == "status":
            status = self.app.device.incoming_status()
            self.app.connection.send_status_message(status)
