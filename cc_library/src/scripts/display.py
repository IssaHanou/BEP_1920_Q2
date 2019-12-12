import os

from cc_library.src.sciler.scclib.app import SccLib
from cc_library.src.sciler.scclib.device import Device

try:
    import RPi.GPIO as GPIO
except (RuntimeError, ModuleNotFoundError):
    from fake_rpi.RPi import GPIO


class Display(Device):
    hint = ""

    def get_status(self):
        status = "{"
        status += "'hint': " + "'" + str(self.hint) + "'"
        status += "}"
        return status

    def perform_instruction(self, contents):
        for action in contents:
            instruction = action.get("instruction")
            if instruction == "hint":
                self.show_hint(action)
            else:
                return False, action
        return True, None

    def test(self):
        self.hint = "test"
        print(self.hint)
        self.scclib.statusChanged()

    def show_hint(self, data):
        self.hint = data.get("value")
        print(self.hint)
        self.scclib.statusChanged()


    def main(self):
        try:
            device = self

            two_up = os.path.abspath(os.path.join(__file__, ".."))
            rel_path = "./display_config.json"
            abs_file_path = os.path.join(two_up, rel_path)
            abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
            config = open(file=abs_file_path)
            self.scclib = SccLib(config=config, device=device)
            self.scclib.start()
        except KeyboardInterrupt:
            self.scclib.logger.log("program was terminated from keyboard input")
        finally:
            GPIO.cleanup()
            self.scclib.logger.log("Cleanly exited Door program")
            self.scclib.logger.close()


if __name__ == "__main__":
    device = Display()
    device.main()
