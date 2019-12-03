# from gpiozero import Button, LED
import RPi.GPIO as GPIO
import time
import math
import Adafruit_ADS1x15

import json
import requests

# def discharge(a_pin, b_pin):
#     GPIO.setup(a_pin, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
#     GPIO.setup(b_pin, GPIO.OUT)
#     GPIO.output(b_pin, GPIO.LOW)
#     time.sleep(0.005)
# def charge_time(a_pin, b_pin):
#     GPIO.setup(b_pin, GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
#     GPIO.setup(a_pin, GPIO.OUT)
#     GPIO.output(a_pin, GPIO.LOW)
#     time.sleep(0.05)
#     count = 0
#     GPIO.output(a_pin, GPIO.HIGH)
#     while not GPIO.input(b_pin):
#         count = count + 1
#     return count
# def analog_read(a_pin, b_pin):
#     discharge(a_pin, b_pin)
#     return charge_time(a_pin, b_pin)
def millis():
    millis = int(round(time.time() * 1000))
    return millis
def within_5(oneval, twoval):
    if abs(oneval-twoval)<5:
        return True
    return False
def get_analog_reading(pin):
    a_pin0 = 24
    a_pin1 = 25
    a_pin2 = 8
    b_pin0 = 7
    b_pin1 = 1
    b_pin2 = 12
    return round(100-(adc.read_adc(pin, gain=1)/266))
def count_true(countlist):
    counter = 0
    for i in range(0,3):
        if (countlist[i]):
            counter = counter+1
    return counter
def correctendposition(potposition, i):
#if i==0:
#    if potposition <2:
#        return True
#    return False
    if within_5(potposition, endPuzzleValue[i]):
        return True
    return False
def correct_position( switchlist):
  counter = 0
  for i in range(0,3):
    if (switchlist[i]==playswitch[i]):
      counter = counter+1
  return counter
def waitForStartupMessage(switches):
  wait = True
  threeneg = False
  while not threeneg:
    potcorrect = 0
    for i in range(0,3):
      potPosition[i] = get_analog_reading(i)
      if potPosition[i]<5:
          potcorrect = potcorrect+1
    print(potPosition, potcorrect)
    if (potcorrect==3):
      threeneg = True
      for i in range(0,3):
          switchPosition[i] = GPIO.input(switches[i])
          if switchPosition[i]:
            threeneg = False
    time.sleep(0.2)
    # ble.println("AT+BLEUARTRX"); TODO check if a start message has been sent by the server. if yes, set threeneg=True, wait=False
    # ble.readline();
    # if (startupString.equals(ble.buffer)){
    #   wait=false;
    #   threeneg = true;
    # }

  pattern = [1,0]
  switchnumber = 0
  times = 0
  while wait:
    newSwitch = GPIO.input(switches[switchnumber])
    # Serial.print(pattern[0]);Serial.print(pattern[1]);Serial.print(switchnumber);Serial.println(newSwitch);
    print(pattern)
    if (newSwitch!=pattern[1]):
      if((pattern[0]==0) and(pattern[1]==1)):
        switchnumber +=1
        pattern[0] = 1
        pattern[1] = 0
      if (switchnumber==3):
        wait=False
      pattern[0] = pattern[1]
      pattern[1] = newSwitch

    time.sleep(0.1)
    # ble.println("AT+BLEUARTRX"); TODO check if a start message has been sent by the server. if yes, set wait=False
    # ble.readline();
    # if (startupString.equals(ble.buffer)){
    #   wait=false;
    # }


adc = Adafruit_ADS1x15.ADS1115()

GPIO.setmode(GPIO.BCM)
blinkTime = [0,0,0]
blinkCount = [999,999,999]
max_blinkcount = 2
winVariable = False
# rood uit, oranje groen aan

correct = [False, False, False]
potPosition = [0,0,0]
prevPotPosition = [0,0,0]
switchPosition = [0,0,0]
goodcount = [0,0,0]
falsecount = [0,0,0]

dramaSwitchPosition = 0
play_puzzle = [False, False, False]

playswitch = [1,1,1]
endPuzzleValue = [3,50,97]
margin = [1,2,100]

startupString = "start"
shutdownString = "shut down"
winMessage = "YOU WIN"
errorSwitchMessage = "wrong switch"


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
greenled = [greenled0, greenled1, greenled2]
redled = [redled0, redled1, redled2]
#waitForStartupMessage(switches)
runTime = millis()

#
# while True:
#     if GPIO.input(switch3):
#         GPIO.output(redled0, GPIO.HIGH)
#         GPIO.output(redled1, GPIO.HIGH)
#         GPIO.output(redled2, GPIO.HIGH)
#     if GPIO.input(switch0):
#         GPIO.output(greenled0, GPIO.HIGH)
#     if GPIO.input(switch1):
#         GPIO.output(greenled1, GPIO.HIGH)
#     if GPIO.input(switch2):
#         GPIO.output(greenled2, GPIO.HIGH)
#     time.sleep(0.2)
#     GPIO.output(redled0, GPIO.LOW)
#     GPIO.output(redled1, GPIO.LOW)
#     GPIO.output(redled2, GPIO.LOW)
#     GPIO.output(greenled0, GPIO.LOW)
#     GPIO.output(greenled1, GPIO.LOW)
#     GPIO.output(greenled2, GPIO.LOW)
#     time.sleep(0.2)
#     print(analog_read(a_pin0, b_pin0),analog_read(a_pin1, b_pin1),analog_read(a_pin2, b_pin2) )

print('hello')
while True:
    if False:#              TODO message from server to reset the scclib
        randommessagevariableNotImportantPleaseDelete = 0
    #     for i in range(0, 3):
    #         GPIO.output(redled[i], GPIO.LOW)
    #         GPIO.output(greenled[i], GPIO.LOW)
    #     waitForStartupMessage(switches)
    #     winVariable = False
    #     dramaSwitchPosition = 0
    #     for i in range(0, 3):
    #         correct[i] = False
    #         play_puzzle[i] = False
    #         potPosition[i] = 0
    #         prevPotPosition[i] = 0
    #         switchPosition[i] = 0
    #         goodcount[i] = 0
    #         potPosition[i] = 0
    #         prevPotPosition[i] = 0
    #         switchPosition[i] = 0
    #         goodcount[i] = 0
    #         blinkTime[i] = 0
    #         blinkCount[i] = 999
    #
    #     for i in range(0, 3):
    #         prevPotPosition[i] = get_analog_reading(i)
    #     runTime = millis()

    # This snippet reads the switches.
    for i in range(0,3):
        switchPosition[i] = GPIO.input(switches[i])
        if playswitch[i] == 0:
            if switchPosition[i] == 0:
                correct[i] = True
            else:
                correct[i] = False
        if switchPosition[i] == 0:
            play_puzzle[i]=False

    # This snippet will brick the board if you turn on the end switch prematurely
    if count_true(correct) < 3 or (correct_position(switchPosition) != 3):
        # if (!digitalRead(dramaSwitch)):
        #     ble.print("AT+BLEUARTTX=")
        #     ble.println(errorSwitchMessage) TODO This part should throw a message to the server that they are trying to cheat, so a feedback sound can be sounded
        while not GPIO.input(switches[3]):
            for i in range(0,3):
                GPIO.output(redled[i], GPIO.HIGH)
                GPIO.output(greenled[i], GPIO.HIGH)
                time.sleep(0.2)
            for i in range(0, 3):
                GPIO.output(redled[i], GPIO.LOW)
                GPIO.output(greenled[i], GPIO.LOW)
                time.sleep(0.2)

    if millis() - runTime > 100:
        runTime = millis()
        dramaSwitchPosition = not GPIO.input(switches[3])
        for i in range(0,3):
            switchPosition[i] = GPIO.input(switches[i])
            potPosition[i] = get_analog_reading(i)
            if not play_puzzle[i] and potPosition[i]!=prevPotPosition[i]:
                if (not correctendposition(potPosition[i], i)) and switchPosition[i] == 1:
                    play_puzzle[i] = True
                    correct[i] = False
                    # GPIO.output(redled[i], GPIO.HIGH)
                    # GPIO.output(greenled[i], GPIO.HIGH)
            # Serial.print(potPosition[0]);
            # Serial.print(" ");
            # Serial.print(potPosition[1]);
            # Serial.print(" ");
            # Serial.print(potPosition[2]);
            # Serial.print(" ");
            # Serial.print(prevPotPosition[0]);
            # Serial.print(" ");
            # Serial.println(play_puzzle[0]);

            # This is the actual puzzle
        print(potPosition, prevPotPosition, play_puzzle, switchPosition, correct, 'hello')
        for i in range(0,3):
            if play_puzzle[i]:
                # Read out the potmeters and see whether they are in the right spot
                potPosition[i] = get_analog_reading(i)
                if correctendposition(potPosition[i], i):
                    prevPotPosition[i] = potPosition[i]
                    goodcount[i] = goodcount[i] + 1
                    falsecount[i] = 0
                else:
                    if within_5(potPosition[i], prevPotPosition[i]):
                        falsecount[i] = falsecount[i] + 1
                    prevPotPosition[i] = potPosition[i]
                    goodcount[i] = 0
                greenledstatus = GPIO.input(greenled[i]) #toggle led pins
                redledstatus = GPIO.input(greenled[i])
                GPIO.output(greenled[i], not greenledstatus)
                GPIO.output(redled[i], not redledstatus)
                if falsecount[i] > 10:
                    falsecount[i]=0
                    blinkCount[i] = 0
                    blinkTime[i] = millis()
                    GPIO.output(redled[i], GPIO.HIGH)
                    GPIO.output(greenled[i], GPIO.HIGH)
                    play_puzzle[i] = False
                    correct[i]=False

                if goodcount[i] > 10:
                    print('its five')
                    goodcount[i]=0
                    blinkCount[i] = 0
                    blinkTime[i] = millis()
                    GPIO.output(redled[i], GPIO.HIGH)
                    GPIO.output(greenled[i], GPIO.HIGH)
                    play_puzzle[i] = False
                    correct[i]=True
                    # if (endPuzzleValue[i]-margin[i] < get_analog_reading(i)) and (endPuzzleValue[i]+margin[i] > get_analog_reading(i)):
                    #     correct[i] = True

    # This piece of code handles status LED's
    for i in range(0,3):
        if not switchPosition[i]:
            GPIO.output(redled[i], GPIO.LOW)
            GPIO.output(greenled[i], GPIO.LOW)
        else:
            if correct[i] and not play_puzzle[i] and (blinkCount[i] > max_blinkcount-1):
                GPIO.output(redled[i], GPIO.LOW)
                GPIO.output(greenled[i], GPIO.HIGH)
            if not correct[i] and not play_puzzle[i] and (blinkCount[i] > max_blinkcount-1):
                GPIO.output(redled[i], GPIO.HIGH)
                GPIO.output(greenled[i], GPIO.LOW)
                # if millis()-blinkTime[i] > 40:
                #     blinkTime[i] = millis()
                #     greenledstatus = GPIO.input(greenled[i])
                #     redledstatus = GPIO.input(greenled[i])
                #     GPIO.output(greenled[i], not greenledstatus)
                #     GPIO.output(redled[i], not redledstatus)
                # GPIO.output(redled[i], GPIO.HIGH)
                # GPIO.output(greenled[i], GPIO.LOW)
            if not play_puzzle[i] and (blinkCount[i] < max_blinkcount):
                if millis()-blinkTime[i] > 40:
                    blinkTime[i] = millis()
                    greenledstatus = GPIO.input(greenled[i])
                    redledstatus = GPIO.input(greenled[i])
                    GPIO.output(greenled[i], not greenledstatus)
                    GPIO.output(redled[i], not redledstatus)
                    blinkCount[i] = blinkCount[i] +1

    # This is the final captive blink loop
    if (count_true(correct) == 3) and (correct_position(switchPosition) == 3) and dramaSwitchPosition:
        # ble.print("AT+BLEUARTTX="); TODO Send to the server a message to indicate victory
        # ble.println(winMessage);
        url = "http://192.168.1.70/api/action/"
        json = {
            "room": "8291cadd-cdca-4ad3-b498-e24686bdcb52",
            "action": "CONTROL_DEVICE",
            "content": {
                "scclib": "SLIDER",
                "solved": True
            }
        }
        r = requests.post(url, json=json)

        winVariable = True
        while winVariable:
            for j in range(0,3):
                for i in range(0,3):
                    if (j == 0):
                        GPIO.output(greenled[i], GPIO.HIGH)
                        time.sleep(0.02)
                        GPIO.output(greenled[i], GPIO.LOW)
                        time.sleep(0.02)
                    else:
                        GPIO.output(redled[2-i], GPIO.HIGH)
                        time.sleep(0.02)
                        GPIO.output(redled[2-i], GPIO.LOW)
                        time.sleep(0.02)
            if  GPIO.input(switches[3]):
                for i in range(0, 3):
                    GPIO.output(redled[i], GPIO.LOW)
                    GPIO.output(greenled[i], GPIO.LOW)
                waitForStartupMessage(switches)
                winVariable = False
                dramaSwitchPosition = 0
                for i in range(0, 3):
                    correct[i] = False
                    play_puzzle[i] = False
                    potPosition[i] = 0
                    prevPotPosition[i] = 0
                    switchPosition[i] = 0
                    goodcount[i] = 0
                    potPosition[i] = 0
                    prevPotPosition[i] = 0
                    switchPosition[i] = 0
                    goodcount[i] = 0
                    blinkTime[i] = 0
                    blinkCount[i] = 999

                for i in range(0, 3):
                    prevPotPosition[i] = get_analog_reading(i)
                runTime = millis()
