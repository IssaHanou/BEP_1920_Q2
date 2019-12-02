import time

try:
    import RPi.GPIO as GPIO
except (RuntimeError, ModuleNotFoundError):
    from fake_rpi.RPi import GPIO


def status(channel):
    print("switch0: ", GPIO.input(switch0))
    print("switch1: ", GPIO.input(switch1))
    print("switch2: ", GPIO.input(switch2))
    print("switch3: ", GPIO.input(switch3))
    print("\n")


try:
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

    a_pin0 = 24
    a_pin1 = 25
    a_pin2 = 8
    b_pin0 = 7
    b_pin1 = 1
    b_pin2 = 12

    GPIO.setup(switch0, GPIO.IN, pull_up_down=GPIO.PUD_UP)
    GPIO.setup(switch1, GPIO.IN, pull_up_down=GPIO.PUD_UP)
    GPIO.setup(switch2, GPIO.IN, pull_up_down=GPIO.PUD_UP)
    GPIO.setup(switch3, GPIO.IN, pull_up_down=GPIO.PUD_UP)
    switches = [switch0, switch1, switch2, switch3]

    GPIO.setup(redled0, GPIO.OUT)
    GPIO.setup(redled1, GPIO.OUT)
    GPIO.setup(redled2, GPIO.OUT)
    GPIO.setup(greenled0, GPIO.OUT)
    GPIO.setup(greenled1, GPIO.OUT)
    GPIO.setup(greenled2, GPIO.OUT)
    greenleds = [greenled0, greenled1, greenled2]
    redleds = [redled0, redled1, redled2]

    GPIO.add_event_detect(switch0, GPIO.BOTH, callback=status)
    GPIO.add_event_detect(switch1, GPIO.BOTH, callback=status)
    GPIO.add_event_detect(switch2, GPIO.BOTH, callback=status)
    GPIO.add_event_detect(switch3, GPIO.BOTH, callback=status)

    while True:
        time.sleep(10 / 1000)

except KeyboardInterrupt:
    print("Interrupted!")

finally:
    GPIO.cleanup()  # This ensures a clean exit
    print("Clean exit ensured!")
