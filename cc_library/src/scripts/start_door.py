import os
import time

from src.device.app import App

try:
    import RPi.GPIO as GPIO
except (RuntimeError, ModuleNotFoundError):
    from fake_rpi.RPi import GPIO


class Door:

    def __init__(self):
        self.door = 17
        GPIO.setup(self.door, GPIO.OUT)

    def incoming_instruction(self, data):
        """
        Set here the mapping from messages to methods.
        """
        print(data)
        instruction = data.get("instruction")
        if instruction == "test":
            self.test()
        elif instruction == "turn off":
            self.turn_off(data)
        elif instruction == "turn on":
            self.turn_on(data),

    def turn_off(self, data):
        #door = getattr(self, data.get("door"))
        GPIO.output(self.door, GPIO.HIGH)

    def turn_on(self, data):
        #door = getattr(self, data.get("door"))
        GPIO.output(self.door, GPIO.LOW)

    def test(self):
        for i in range(0, 2):
            GPIO.output(self.door, GPIO.HIGH)
            time.sleep(5)
            GPIO.output(self.door, GPIO.LOW)
            time.sleep(5)
        GPIO.output(self.door, GPIO.HIGH)

if __name__ == "__main__":
    # This is the main script of the Raspberry Pi
    GPIO.setmode(GPIO.BCM)

    device1 = Door()

    two_up = os.path.abspath(os.path.join(__file__, ".."))
    rel_path = "./door_config.json"
    abs_file_path = os.path.join(two_up, rel_path)
    abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
    config = open(file=abs_file_path)
    app = App(config=config, device=device1)

    app.start()

