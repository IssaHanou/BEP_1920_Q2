import os
import time

from cc_library.src.sciler.scclib.app import SccLib
from cc_library.src.sciler.scclib.device import Device

try:
    import RPi.GPIO as GPIO
except (RuntimeError, ModuleNotFoundError):
    from fake_rpi.RPi import GPIO

scclib = None


class Door(Device):
    """
    Define pin numbers to which units are connected on Pi.
    """

    GPIO.setmode(GPIO.BCM)
    door = 17
    GPIO.setup(door, GPIO.OUT)

    status = False

    def get_status(self):
        """
        Return status of the door.
        """
        status = "{"
        status += "'door': "
        if self.status:
            status += "'open'"
        else:
            status += "'closed'"
        status += "}"
        return status

    def perform_instruction(self, contents):
        """
        Set here the mapping from messages to methods.
        Should return warning when illegal instruction was sent
        or instruction could not be performed.
        """
        for action in contents:
            instruction = action.get("instruction")
            if instruction == "door":
                if action.get("value"):
                    self.turn_off()
                else:
                    self.turn_on()
            else:
                return True
        return False

    def test(self):
        for i in range(0, 2):
            self.turn_on()
            time.sleep(2)
            self.turn_off()
            time.sleep(2)
        self.turn_on()

    def turn_off(self):
        GPIO.output(self.door, GPIO.HIGH)
        self.status = False
        scclib.statusChanged()

    def turn_on(self):
        GPIO.output(self.door, GPIO.LOW)
        self.status = True
        scclib.statusChanged()


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
    scclib.logger.log("program was terminated from keyboard input")
finally:
    GPIO.cleanup()
    scclib.logger.log("Cleanly exited Door program")
    scclib.logger.close()
