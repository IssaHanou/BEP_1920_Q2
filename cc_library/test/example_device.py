from cc_library.src.device.app import App

try:
    import RPi.GPIO as GPIO
except (RuntimeError, ModuleNotFoundError):
    import fake_rpi.RPi as GPIO
import time

GPIO.setmode(GPIO.BCM)

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
switches = [switch0, switch1, switch2, switch3]
switchPosition = [0, 0, 0]
for i in range(0, 3):
    switchPosition[i] = GPIO.input(switches[i])

GPIO.setup(redled0, GPIO.OUT)
GPIO.setup(redled1, GPIO.OUT)
GPIO.setup(redled2, GPIO.OUT)
GPIO.setup(greenled0, GPIO.OUT)
GPIO.setup(greenled1, GPIO.OUT)
GPIO.setup(greenled2, GPIO.OUT)
greenled = [greenled0, greenled1, greenled2]
redled = [redled0, redled1, redled2]


def blink(led, interval):
    GPIO.output(led, GPIO.HIGH)
    time.sleep(interval)
    GPIO.output(led, GPIO.LOW)
    time.sleep(interval)


def turnOff(led):
    GPIO.output(led, GPIO.LOW)


def turnOn(led):
    GPIO.output(led, GPIO.HIGH)


def getSwitches(switchPosition=None):
    for i in range(0, 3):
        switchPosition[i] = GPIO.input(switches[i])


def test():
    for j in range(0, 10):
        for i in range(0, 3):
            GPIO.output(redled[i], GPIO.HIGH)
            GPIO.output(greenled[i], GPIO.HIGH)
            time.sleep(0.2)
        for i in range(0, 3):
            GPIO.output(redled[i], GPIO.LOW)
            GPIO.output(greenled[i], GPIO.LOW)
            time.sleep(0.2)


config = open('example_config.json')
app = App(config)

app.start()
