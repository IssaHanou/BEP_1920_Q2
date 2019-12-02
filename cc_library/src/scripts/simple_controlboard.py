import os
import time

from src.scclib.scclib import SccLib

try:
    import RPi.GPIO as GPIO
except (RuntimeError, ModuleNotFoundError):
    from fake_rpi.RPi import GPIO


class ControlBoard:
    GPIO.setmode(GPIO.BCM)

    redSwitch = 27
    orangeSwitch = 22
    greenSwitch = 18
    mainSwitch = 23
    switches = [redSwitch, orangeSwitch, greenSwitch, mainSwitch]

    redled0 = 9
    redled1 = 15
    redled2 = 17
    greenled0 = 10
    greenled1 = 14
    greenled2 = 4

    a_pin0 = 24
    a_pin1 = 25
    a_pin2 = 8
    b_pin0 = 7
    b_pin1 = 1
    b_pin2 = 12

    GPIO.setup(redSwitch, GPIO.IN, pull_up_down=GPIO.PUD_UP)
    GPIO.setup(orangeSwitch, GPIO.IN, pull_up_down=GPIO.PUD_UP)
    GPIO.setup(greenSwitch, GPIO.IN, pull_up_down=GPIO.PUD_UP)
    GPIO.setup(mainSwitch, GPIO.IN, pull_up_down=GPIO.PUD_UP)

    GPIO.setup(redled0, GPIO.OUT)
    GPIO.setup(redled1, GPIO.OUT)
    GPIO.setup(redled2, GPIO.OUT)
    GPIO.setup(greenled0, GPIO.OUT)
    GPIO.setup(greenled1, GPIO.OUT)
    GPIO.setup(greenled2, GPIO.OUT)

    def get_status(self):
        status = "{"
        status += ("'redSwitch': " + str(GPIO.input(self.redSwitch)) + ",")
        status += ("'orangeSwitch': " + str(GPIO.input(self.orangeSwitch)) + ",")
        status += ("'greenSwitch': " + str(GPIO.input(self.greenSwitch)) + ",")
        status += ("'mainSwitch': " + str(GPIO.input(self.mainSwitch)))
        status += "}"
        return status

    def perform_instruction(self, contents):
        """
        Set here the mapping from messages to methods.
        """
        instruction = contents.get("instruction")
        if instruction == "test":
            self.test()
        elif instruction == "blink":
            self.blink(contents)
        elif instruction == "turnOff":
            self.turn_off(contents)
        elif instruction == "turnOn":
            self.turn_on(contents)

    def blink(self, data):
        led = getattr(self, data.get("led"))
        interval = data.get("interval")
        GPIO.output(led, GPIO.HIGH)
        time.sleep(interval)
        GPIO.output(led, GPIO.LOW)
        time.sleep(interval)

    def turn_off(self, data):
        led = getattr(self, data.get("led"))
        GPIO.output(led, GPIO.LOW)

    def turn_on(self, data):
        led = getattr(self, data.get("led"))
        GPIO.output(led, GPIO.HIGH)

    def test(self):
        for j in range(0, 3):
            for i in range(0, 3):
                GPIO.output(self.redled[i], GPIO.HIGH)
                GPIO.output(self.greenled[i], GPIO.HIGH)
                time.sleep(0.2)
            for i in range(0, 3):
                GPIO.output(self.redled[i], GPIO.LOW)
                GPIO.output(self.greenled[i], GPIO.LOW)
                time.sleep(0.2)


try:
    device = ControlBoard()

    two_up = os.path.abspath(os.path.join(__file__, ".."))
    rel_path = "./controlboard_config.json"
    abs_file_path = os.path.join(two_up, rel_path)
    abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
    config = open(file=abs_file_path)
    scclib = SccLib(config, device)

    GPIO.add_event_detect(device.redSwitch, GPIO.BOTH, callback=scclib.statusChanged)
    GPIO.add_event_detect(device.orangeSwitch, GPIO.BOTH, callback=scclib.statusChanged)
    GPIO.add_event_detect(device.greenSwitch, GPIO.BOTH, callback=scclib.statusChanged)
    GPIO.add_event_detect(device.mainSwitch, GPIO.BOTH, callback=scclib.statusChanged)

    scclib.start()
except KeyboardInterrupt:
    print("Interrupted!")

finally:
    GPIO.cleanup()  # This ensures a clean exit
    print("Clean exit ensured!")
