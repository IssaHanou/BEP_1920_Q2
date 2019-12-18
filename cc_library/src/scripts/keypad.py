import os
from cc_library.src.sciler.scclib.device import Device


try:
    from evdev import InputDevice, categorize, ecodes
except ImportError:
    print("Error import EVDEV. Module Room_door_handler can't be used")

import asyncio


class Keypad(Device):
    currentValue = ""

    def get_status(self):
        """
        Returns status of all custom components, in json format.
        """

        return {"code": self.currentValue}

    def perform_instruction(self, contents):
        """
        Defines how instructions are handled,
        for all instructions defined in output of device in config.
        :param contents: contains instruction tag and calls the appropriate functions.
        :return boolean: True if instruction was valid and False if illegal instruction
        was sent or error occurred such that instruction could not be performed.
        Returns tuple, with boolean and None if True and the failed action if false.
        """

    def test(self):
        """
        Defines test sequence for device.
        """

    def reset(self):
        """
        Defines a reset sequence for device.
        """

    def __init__(self):
        two_up = os.path.abspath(os.path.join(__file__, ".."))
        rel_path = "./keypad.json"
        abs_file_path = os.path.join(two_up, rel_path)
        abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
        config = open(file=abs_file_path)
        super().__init__(config)

    def main(self):
        """
        This method should be overridden in the subclass.
        It should initialize the SccLib class with config file and device class.
        It should also add event listeners to GPIO for all input components.
        """

        def async_loop():
            card_reader = NumpadReader(
                result_handler=self.reader_handle_result,
                status_change=self.reader_status_changed,
                log=self.log,
                usb_path="/dev/input/event0",
            )
            loop = asyncio.get_event_loop()

            loop.run_until_complete(asyncio.gather(card_reader.start()))

        self.start(loop=async_loop)

    def reader_handle_result(self, result):
        self.log("Submitted code: {}".format(result))
        self.currentValue = result
        self.status_changed()

    def reader_status_changed(self, result):
        self.log("Current status: {}".format(result))
        self.currentValue = result
        self.status_changed()


class NumpadReader:
    def __init__(
        self, result_handler, status_change, usb_path="/dev/input/event0", log=print
    ):
        self.word = ""
        self.device = InputDevice(usb_path)
        self.count = 0
        self.handler = result_handler
        self.statusChange = status_change
        self.log = log

    async def start(self):
        self.log("Started listening for input")
        async for event in self.device.async_read_loop():
            if event.type == ecodes.EV_KEY:
                data = categorize(event)
                if data.keystate == 1:
                    value = data.scancode
                    self.process_value(value)

    def process_value(self, value):
        self.count = self.count + 1
        number = self.get_value(value)

        if not self.check_enter(value):
            if number is None:
                return

            self.word = self.word + str(number)
            self.statusChange(self.word)
            return

        # Could be something for the config

        self.handler(self.word)

        self.count = 0
        self.word = ""

    def check_enter(self, c):
        return c == 28 or c == 96  # KEY_ENTER or KEY_KPENTER

    def get_value(self, value):
        mapping = {
            102: 7,
            71: 7,
            103: 8,
            72: 8,
            104: 9,
            73: 9,
            75: 4,
            105: 4,
            76: 5,
            106: 6,
            77: 6,
            107: 1,
            79: 1,
            108: 2,
            80: 2,
            109: 3,
            81: 3,
            110: 0,
            82: 0,
        }
        return mapping.get(value, None)


if __name__ == "__main__":
    device = Keypad()
    device.main()
