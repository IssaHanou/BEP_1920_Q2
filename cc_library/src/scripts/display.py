import os

from cc_library.src.sciler.scclib.app import SccLib
from cc_library.src.sciler.scclib.device import Device

try:
    import RPi.GPIO as GPIO
except (RuntimeError, ModuleNotFoundError):
    from fake_rpi.RPi import GPIO

scclib = None


class Display(Device):
    hint = ""

    def get_status(self):
        status = "{"
        status += "'hint': " + "'" + str(self.hint) + "'"
        status += "}"
        return status

    def perform_instruction(self, contents):
        instruction = contents.get("instruction")
        if instruction == "test":
            self.test()
        elif instruction == "hint":
            self.show_hint(contents)
        else:
            return True
        return None

    def test(self):
        self.hint = "test"
        print(self.hint)
        scclib.statusChanged()

    def show_hint(self, data):
        self.hint = data.get("hint")
        print(self.hint)
        scclib.statusChanged()


try:
    device = Display()

    two_up = os.path.abspath(os.path.join(__file__, ".."))
    rel_path = "./display_config.json"
    abs_file_path = os.path.join(two_up, rel_path)
    abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
    config = open(file=abs_file_path)
    scclib = SccLib(config=config, device=device)
    scclib.start()
except KeyboardInterrupt:
    scclib.logger.log("program was terminated from keyboard input")
finally:
    GPIO.cleanup()
    scclib.logger.log("Cleanly exited Door program")
    scclib.logger.close()
