import os
import time

from cc_library.src.sciler.scclib.app import SccLib
from cc_library.src.sciler.scclib.device import Device
import Adafruit_ADS1x15

try:
    import RPi.GPIO as GPIO
except (RuntimeError, ModuleNotFoundError):
    from fake_rpi.RPi import GPIO


class ControlBoard(Device):
    GPIO.setmode(GPIO.BCM)
    """
    Define pin numbers to which units are connected on Pi.
    """
    redSwitch = 27
    orangeSwitch = 22
    greenSwitch = 18
    mainSwitch = 23
    switches = [redSwitch, orangeSwitch, greenSwitch, mainSwitch]

    redLED0 = 9
    redLED1 = 15
    redLED2 = 17
    redLEDs = [redLED0, redLED1, redLED2]
    greenLED0 = 10
    greenLED1 = 14
    greenLED2 = 4
    greenLEDs = [greenLED0, greenLED1, greenLED2]

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

    GPIO.setup(redLED0, GPIO.OUT)
    GPIO.setup(redLED1, GPIO.OUT)
    GPIO.setup(redLED2, GPIO.OUT)
    GPIO.setup(greenLED0, GPIO.OUT)
    GPIO.setup(greenLED1, GPIO.OUT)
    GPIO.setup(greenLED2, GPIO.OUT)

    adc = Adafruit_ADS1x15.ADS1115()

    GPIO.setup(a_pin0, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
    GPIO.setup(a_pin1, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
    GPIO.setup(a_pin2, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
    GPIO.setup(b_pin0, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
    GPIO.setup(b_pin1, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
    GPIO.setup(b_pin2, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)

    def get_sliders_analog_reading(self):
        positions = [0, 0, 0]
        for channel in range(0, 3):
            positions[channel] = round(100 - (self.adc.read_adc(channel) / 266))
        return positions

    def get_leds_status(self):
        reds = [0, 0, 0]
        greens = [0, 0, 0]
        for i in range(0, 3):
            reds[i] = GPIO.input(self.redLEDs[i])
            greens[i] = GPIO.input(self.greenLEDs[i])
        return "'redLEDs': " + str(reds) + ",'greenLEDs': " + str(greens)

    def get_status(self):
        """
        Return status of switches, LEDs and sliders of device.
        """
        status = "{"
        status += "'redSwitch': " + str(GPIO.input(self.redSwitch)) + ","
        status += "'orangeSwitch': " + str(GPIO.input(self.orangeSwitch)) + ","
        status += "'greenSwitch': " + str(GPIO.input(self.greenSwitch)) + ","
        status += "'mainSwitch': " + str(GPIO.input(self.mainSwitch)) + ","
        status += self.get_leds_status() + ","
        status += "'sliders': " + str(self.get_sliders_analog_reading())
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

            if instruction == "blink":
                self.blink(action.get("component_id"), action.get("value"))
            elif instruction == "turnOnOff":
                self.turn_on_off(action.get("component_id"), action.get("value"))
            else:
                return True
        return False

    def blink(self, component, args):
        led = getattr(self, component)
        time.sleep(args[1])  # delay
        interval = args[0]
        GPIO.output(led, GPIO.HIGH)
        time.sleep(interval)
        GPIO.output(led, GPIO.LOW)
        time.sleep(interval)

    def turn_on_off(self, component, arg):
        led = getattr(self, component)
        if arg:
            GPIO.output(led, GPIO.LOW)
        else:
            GPIO.output(led, GPIO.HIGH)

    def test(self):
        for j in range(0, 3):
            for i in range(0, 3):
                GPIO.output(self.redLEDs[i], GPIO.HIGH)
                GPIO.output(self.greenLEDs[i], GPIO.HIGH)
                time.sleep(0.2)
            for i in range(0, 3):
                GPIO.output(self.redLEDs[i], GPIO.LOW)
                GPIO.output(self.greenLEDs[i], GPIO.LOW)
                time.sleep(0.2)


scclib = None
try:
    device = ControlBoard()

    two_up = os.path.abspath(os.path.join(__file__, ".."))
    rel_path = "./controlboard_config.json"
    abs_file_path = os.path.join(two_up, rel_path)
    abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
    config = open(file=abs_file_path)
    scclib = SccLib(config, device)

    GPIO.add_event_detect(
        device.redSwitch, GPIO.BOTH, callback=scclib.statusChangedOnChannel, bouncetime=100
    )
    GPIO.add_event_detect(
        device.orangeSwitch, GPIO.BOTH, callback=scclib.statusChangedOnChannel, bouncetime=100
    )
    GPIO.add_event_detect(
        device.greenSwitch, GPIO.BOTH, callback=scclib.statusChangedOnChannel, bouncetime=100
    )
    GPIO.add_event_detect(
        device.mainSwitch, GPIO.BOTH, callback=scclib.statusChangedOnChannel, bouncetime=100
    )
    GPIO.add_event_detect(device.a_pin0, GPIO.BOTH, callback=scclib.status_changed)
    GPIO.add_event_detect(device.a_pin1, GPIO.BOTH, callback=scclib.status_changed)
    GPIO.add_event_detect(device.a_pin2, GPIO.BOTH, callback=scclib.status_changed)
    GPIO.add_event_detect(device.b_pin0, GPIO.BOTH, callback=scclib.status_changed)
    GPIO.add_event_detect(device.b_pin1, GPIO.BOTH, callback=scclib.status_changed)
    GPIO.add_event_detect(device.b_pin2, GPIO.BOTH, callback=scclib.status_changed)

    scclib.start()
except KeyboardInterrupt:
    scclib.logger.log("program was terminated from keyboard input")
finally:
    GPIO.cleanup()
    scclib.logger.log("Cleanly exited ControlBoard program")
    scclib.logger.close()
