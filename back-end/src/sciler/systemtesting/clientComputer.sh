#!/bin/bash

echo "raspberry" | sudo -S apt update
echo "raspberry" | sudo -S apt-get install python3.5
echo "raspberry" | echo "Y" | sudo -S apt install python3-pip
echo "raspberry" | sudo -S -H pip3 install scclib
echo "raspberry" | sudo -S -H pip3 install paho-mqtt
python3 latency.py &
echo "Client computer running"