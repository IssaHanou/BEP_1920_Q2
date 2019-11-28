import os
import time
from src.device.app import App

try:
    import RPi.GPIO as GPIO
except (RuntimeError, ModuleNotFoundError):
    from fake_rpi.RPi import GPIO


class ControlBoard:
    """
    This class is written by the programmer of the Raspberry Pi.
    It should contain a incoming_instruction(data)
    with a mapping of messages to methods.
    """

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

        self.switch0 = switch0
        self.switch1 = switch1
        self.switch2 = switch2
        self.switch3 = switch3

        self.redled0 = redled0
        self.redled1 = redled1
        self.redled2 = redled2
        self.greenled0 = greenled0
        self.greenled1 = greenled1
        self.greenled2 = greenled2

    def incoming_instruction(self, data):
        """
        Set here the mapping from messages to methods.
        """
        print(data)
        instruction = data.get("instruction")
        if instruction == "test":
            self.test()
        elif instruction == "blink":
            self.blink(data)
        elif instruction == "turnOff":
            self.turn_off(data)
        elif instruction == "turnOn":
            self.turn_on(data),

    def incoming_status(self):
        return self.get_switch()

    def blink(self, data):
        led = getattr(self, data.get("led"))
        interval = data.get("interval")
        print("blinking: ", led, interval)
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

    def get_switch(self, name, switch):
        print({name: str(GPIO.input(switch))})
        return str({name: str(GPIO.input(switch))})

    def test(self):
        print("debug2")
        for j in range(0, 3):
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
    # This is the main script of the Raspberry Pi
    GPIO.setmode(GPIO.BCM)

    device1 = ControlBoard()

    two_up = os.path.abspath(os.path.join(__file__, ".."))
    rel_path = "./controlboard_config.json"
    abs_file_path = os.path.join(two_up, rel_path)
    abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
    config = open(file=abs_file_path)
    app = App(config=config, device=device1)

    def get_status_switch(device, name, switch):
        """
        get status calls send_message to send the status of
        a switch to the MQTT broker.
        """
        print("get status", name, switch)
        app.send_status_message(device.get_switch(name, switch))

    GPIO.add_event_detect(
        device1.switch0,
        GPIO.BOTH,
        callback=lambda *a: get_status_switch(device1, "switch0", device1.switch0),
        bouncetime=100,
    )
    GPIO.add_event_detect(
        device1.switch1,
        GPIO.BOTH,
        callback=lambda *a: get_status_switch(device1, "switch1", device1.switch1),
        bouncetime=100,
    )
    GPIO.add_event_detect(
        device1.switch2,
        GPIO.BOTH,
        callback=lambda *a: get_status_switch(device1, "switch2", device1.switch2),
        bouncetime=100,
    )
    GPIO.add_event_detect(
        device1.switch3,
        GPIO.BOTH,
        callback=lambda *a: get_status_switch(device1, "switch3", device1.switch3),
        bouncetime=100,
    )

    # Start connection
    app.start()
