import os
import time

from scclib.device import Device

try:
    import RPi.GPIO as GPIO
except (RuntimeError, ModuleNotFoundError):
    from fake_rpi.RPi import GPIO


class Door(Device):
    def __init__(self):
        """
        Define pin numbers to which units are connected on Pi.
        """
        two_up = os.path.abspath(os.path.join(__file__, ".."))
        rel_path = "./door_config.json"
        abs_file_path = os.path.join(two_up, rel_path)
        abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
        config = open(file=abs_file_path)
        super().__init__(config)
        GPIO.setmode(GPIO.BCM)
        self.door = 17
        GPIO.setup(self.door, GPIO.OUT)
        self.status = False

    def get_status(self):
        """
        Return status of the door.
        """
        status_str = "closed"
        if self.status:
            status_str = "open"
        return {"door": status_str}

    def perform_instruction(self, action):
        """
        Set here the mapping from messages to methods.
        Should return warning when illegal instruction was sent
        or instruction could not be performed.
        """
        instruction = action.get("instruction")
        if instruction == "open":
            if action.get("value"):
                self.turn_off()
            else:
                self.turn_on()
        else:
            return False, action
        return True, None

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
        self.status_changed()

    def turn_on(self):
        GPIO.output(self.door, GPIO.LOW)
        self.status = True
        self.status_changed()

    def reset(self):
        self.turn_off()

    def main(self):
        self.start(stop=GPIO.cleanup)


if __name__ == "__main__":
    device = Door()
    device.main()
