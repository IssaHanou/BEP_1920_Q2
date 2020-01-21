import os
import time

from scclib.device import Device
from Adafruit_ADS1x15 import ADS1115 as ADC

try:
    import RPi.GPIO as GPIO
except (RuntimeError, ModuleNotFoundError):
    from fake_rpi.RPi import GPIO


class ControlBoard(Device):
    def __init__(self):

        two_up = os.path.abspath(os.path.join(__file__, ".."))
        rel_path = "./controlboard_config.json"
        abs_file_path = os.path.join(two_up, rel_path)
        abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
        config = open(file=abs_file_path)

        super().__init__(config)
        self.adc = ADC()
        GPIO.setmode(GPIO.BCM)
        """
        Define pin numbers to which units are connected on Pi.
        """
        self.redSwitch = 27
        self.orangeSwitch = 22
        self.greenSwitch = 18
        self.mainSwitch = 23
        self.switches = [
            self.redSwitch,
            self.orangeSwitch,
            self.greenSwitch,
            self.mainSwitch,
        ]

        self.redLight1 = 9
        self.redLight2 = 15
        self.redLight3 = 17
        self.greenLight1 = 10
        self.greenLight2 = 14
        self.greenLight3 = 4

        self.redLEDs = [self.redLight1, self.redLight2, self.redLight3]
        self.greenLEDs = [self.greenLight1, self.greenLight2, self.greenLight3]

        self.a_pin0 = 24
        self.a_pin1 = 25
        self.a_pin2 = 8
        self.b_pin0 = 7
        self.b_pin1 = 1
        self.b_pin2 = 12

        GPIO.setup(self.redSwitch, GPIO.IN, pull_up_down=GPIO.PUD_UP)
        GPIO.setup(self.orangeSwitch, GPIO.IN, pull_up_down=GPIO.PUD_UP)
        GPIO.setup(self.greenSwitch, GPIO.IN, pull_up_down=GPIO.PUD_UP)
        GPIO.setup(self.mainSwitch, GPIO.IN, pull_up_down=GPIO.PUD_UP)

        GPIO.setup(self.redLight1, GPIO.OUT)
        GPIO.setup(self.redLight2, GPIO.OUT)
        GPIO.setup(self.redLight3, GPIO.OUT)
        GPIO.setup(self.greenLight1, GPIO.OUT)
        GPIO.setup(self.greenLight2, GPIO.OUT)
        GPIO.setup(self.greenLight3, GPIO.OUT)

        GPIO.setup(self.a_pin0, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
        GPIO.setup(self.a_pin1, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
        GPIO.setup(self.a_pin2, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
        GPIO.setup(self.b_pin0, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
        GPIO.setup(self.b_pin1, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
        GPIO.setup(self.b_pin2, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)

    def get_sliders_analog_reading(self):
        positions = [0, 0, 0]
        for channel in range(0, 3):
            positions[channel] = round(100 - (self.adc.read_adc(channel) / 266))
        return positions

    def status_binair_to_bool(self, binair):
        if binair == 0:
            return False
        else:
            return True

    def status_binair_to_sting(self, binair):
        if binair == 0:
            return "uit"
        else:
            return "aan"

    def get_status(self):
        """
        Return status of switches, LEDs and sliders of device.
        """
        return {
            "redSwitch": self.status_binair_to_bool(GPIO.input(self.redSwitch)),
            "orangeSwitch": self.status_binair_to_bool(GPIO.input(self.orangeSwitch)),
            "greenSwitch": self.status_binair_to_bool(GPIO.input(self.greenSwitch)),
            "mainSwitch": self.status_binair_to_bool(GPIO.input(self.mainSwitch)),
            "greenLight1": self.status_binair_to_sting(GPIO.input(self.greenLight1)),
            "greenLight2": self.status_binair_to_sting(GPIO.input(self.greenLight2)),
            "greenLight3": self.status_binair_to_sting(GPIO.input(self.greenLight3)),
            "redLight1": self.status_binair_to_sting(GPIO.input(self.redLight1)),
            "redLight2": self.status_binair_to_sting(GPIO.input(self.redLight2)),
            "redLight3": self.status_binair_to_sting(GPIO.input(self.redLight3)),
            "slider1": self.get_sliders_analog_reading()[0],
            "slider2": self.get_sliders_analog_reading()[1],
            "slider3": self.get_sliders_analog_reading()[2],
        }

    def perform_instruction(self, action):
        """
        Set here the mapping from messages to methods.
        Should return warning when illegal instruction was sent
        or instruction could not be performed.
        """
        instruction = action.get("instruction")
        if instruction == "blink":
            self.blink(action.get("component_id"), action.get("value"))
        elif instruction == "turnOnOff":
            self.turn_on_off(action.get("component_id"), action.get("value"))
        else:
            return False, action

        return True, None

    def blink(self, component, args):
        led = getattr(self, component)

        time.sleep(args[1])  # delay
        interval = args[0]
        GPIO.output(led, GPIO.HIGH)
        time.sleep(interval)
        GPIO.output(led, GPIO.LOW)
        time.sleep(interval)
        GPIO.output(led, GPIO.HIGH)
        time.sleep(interval)
        GPIO.output(led, GPIO.LOW)
        time.sleep(interval)

    def turn_on_off(self, component, arg):
        led = getattr(self, component)
        if arg:
            GPIO.output(led, GPIO.HIGH)
        else:
            GPIO.output(led, GPIO.LOW)

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

    def __status_update(self, channel):
        self.status_changed()

    def setup_events(self):

        GPIO.add_event_detect(
            device.redSwitch, GPIO.BOTH, callback=self.__status_update, bouncetime=100,
        )
        GPIO.add_event_detect(
            device.orangeSwitch,
            GPIO.BOTH,
            callback=self.__status_update,
            bouncetime=100,
        )
        GPIO.add_event_detect(
            device.greenSwitch,
            GPIO.BOTH,
            callback=self.__status_update,
            bouncetime=100,
        )
        GPIO.add_event_detect(
            device.mainSwitch, GPIO.BOTH, callback=self.__status_update, bouncetime=100,
        )
        GPIO.add_event_detect(
            device.a_pin0, GPIO.BOTH, callback=self.__status_update,
        )
        GPIO.add_event_detect(
            device.a_pin1, GPIO.BOTH, callback=self.__status_update,
        )
        GPIO.add_event_detect(
            device.a_pin2, GPIO.BOTH, callback=self.__status_update,
        )
        GPIO.add_event_detect(
            device.b_pin0, GPIO.BOTH, callback=self.__status_update,
        )
        GPIO.add_event_detect(
            device.b_pin1, GPIO.BOTH, callback=self.__status_update,
        )
        GPIO.add_event_detect(
            device.b_pin2, GPIO.BOTH, callback=self.__status_update,
        )

    def loop(self):
        previous = self.get_sliders_analog_reading()
        while True:
            positions = self.get_sliders_analog_reading()
            if previous != positions:
                self.status_changed()
                previous = positions

    def reset(self):
        for i in range(0, 3):
            GPIO.output(self.redLEDs[i], GPIO.LOW)
            GPIO.output(self.greenLEDs[i], GPIO.LOW)

    def main(self):
        self.setup_events()
        self.start(loop=self.loop, stop=GPIO.cleanup)


if __name__ == "__main__":
    device = ControlBoard()
    device.main()
