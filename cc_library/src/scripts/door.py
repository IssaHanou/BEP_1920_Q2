import os
import time

from cc_library.src.sciler.scclib.app import SccLib, on_python_log
from cc_library.src.sciler.scclib.device import Device

try:
    import RPi.GPIO as GPIO
except (RuntimeError, ModuleNotFoundError):
    from fake_rpi.RPi import GPIO


class Door(Device):
    """
    Define pin numbers to which units are connected on Pi.
    """
    GPIO.setmode(GPIO.BCM)
    door = 17
    GPIO.setup(door, GPIO.OUT)

    def get_status(self):
        """
        Return status of the door.
        """
        status = "{"
        status += "'door': " + str(GPIO.input(self.door))
        status += "}"
        return status

    def perform_instruction(self, contents):
        """
        Set here the mapping from messages to methods.
        """
        instruction = contents.get("instruction")
        if instruction == "test":
            self.test()
        elif instruction == "turn off":
            self.turn_off(contents)
        elif instruction == "turn on":
            self.turn_on(contents),

    def test(self):
        for i in range(0, 2):
            GPIO.output(self.door, GPIO.LOW)
            time.sleep(2)
            GPIO.output(self.door, GPIO.HIGH)
            time.sleep(2)
        GPIO.output(self.door, GPIO.LOW)

    def turn_off(self, data):
        door = getattr(self, data.get("door"))
        GPIO.output(door, GPIO.HIGH)

    def turn_on(self, data):
        door = getattr(self, data.get("door"))
        GPIO.output(door, GPIO.LOW)


try:
    device = Door()

    two_up = os.path.abspath(os.path.join(__file__, ".."))
    rel_path = "./door_config.json"
    abs_file_path = os.path.join(two_up, rel_path)
    abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
    config = open(file=abs_file_path)
    scclib = SccLib(config=config, device=device)
    scclib.start()

except KeyboardInterrupt:
    on_python_log("Interrupted by keyboard")
finally:
    GPIO.cleanup()  # This ensures clean exit
    on_python_log("Cleanly exited Door program")