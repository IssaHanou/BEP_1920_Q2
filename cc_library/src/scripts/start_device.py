print("starting")
import os

from src.device.app import App

try:
    import RPi.GPIO as GPIO
except (RuntimeError, ModuleNotFoundError):
    from fake_rpi.RPi import GPIO
import time


class ControlBoard:
    def __init__(self):
        switch0 = 27
        switch1 = 22
        switch2 = 18
        switch3 = 23

        redled0 = 9
        redled1 = 15
        redled2 = 17
        greenled0 = 10
        greenled1 = 14
        greenled2 = 4

        GPIO.setup(switch0, GPIO.IN, pull_up_down=GPIO.PUD_UP)
        GPIO.setup(switch1, GPIO.IN, pull_up_down=GPIO.PUD_UP)
        GPIO.setup(switch2, GPIO.IN, pull_up_down=GPIO.PUD_UP)
        GPIO.setup(switch3, GPIO.IN, pull_up_down=GPIO.PUD_UP)
        self.switches = [switch0, switch1, switch2, switch3]
        self.switchPosition = [0, 0, 0]
        for i in range(0, 3):
            self.switchPosition[i] = GPIO.input(self.switches[i])

        GPIO.setup(redled0, GPIO.OUT)
        GPIO.setup(redled1, GPIO.OUT)
        GPIO.setup(redled2, GPIO.OUT)
        GPIO.setup(greenled0, GPIO.OUT)
        GPIO.setup(greenled1, GPIO.OUT)
        GPIO.setup(greenled2, GPIO.OUT)
        self.greenled = [greenled0, greenled1, greenled2]
        self.redled = [redled0, redled1, redled2]

    def blink(self, led, interval):
        GPIO.output(self, GPIO.HIGH)
        time.sleep(interval)
        GPIO.output(self, GPIO.LOW)
        time.sleep(interval)

    def turnOff(self, led):
        GPIO.output(led, GPIO.LOW)

    def turnOn(self, led):
        GPIO.output(led, GPIO.HIGH)

    def getSwitches(self, switchPosition=None):
        for i in range(0, 3):
            switchPosition[i] = GPIO.input(self.switches[i])

    def test(self):
        for j in range(0, 10):
            for i in range(0, 3):
                GPIO.output(self.redled[i], GPIO.HIGH)
                GPIO.output(self.greenled[i], GPIO.HIGH)
                time.sleep(0.2)
            for i in range(0, 3):
                GPIO.output(self.redled[i], GPIO.LOW)
                GPIO.output(self.greenled[i], GPIO.LOW)
                time.sleep(0.2)
            print("knipper knipper")


if __name__ == "__main__":
    GPIO.setmode(GPIO.BCM)
    device1 = ControlBoard()

    two_up = os.path.abspath(os.path.join(__file__, ".."))
    rel_path = "./example_config.json"
    abs_file_path = os.path.join(two_up, rel_path)
    abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
    config = open(file=abs_file_path)
    app = App(config=config, device=device1)

    app.start()
