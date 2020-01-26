#!/bin/bash

sudo apt-get update
sudo apt-get install -y python3.5
sudo apt install -y python3-pip
sudo -H pip3 install scclib
sudo -H pip3 install paho-mqtt
python3 latency.py &
echo "Client computer running"