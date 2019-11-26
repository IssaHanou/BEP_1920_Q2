import json


# from src.scripts import start_device


class BaseHandler:
    def __init__(self):
        self.status = None

    def on_change(self, status):
        self.status = status

    def on_test(self):
        self.app.device.test()
        print("testing")

    def handle(self, message):
        message = message.payload.decode("utf-8")
        message = json.loads(message)
        message_type = message.get("messageInstructionTest").get("type")
        if message_type == "instruction":
            message_instruction = (
                message.get("messageInstructionTest").get("contents").get("instruction")
            )
            if message_instruction == "test":
                BaseHandler.on_test(self)
