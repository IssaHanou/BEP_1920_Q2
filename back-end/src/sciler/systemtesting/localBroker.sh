#!/bin/bash

echo raspberry | sudo -S apt update
echo raspberry | sudo -S apt install -y mosquitto mosquitto-clients
echo raspberry | sudo -S service mosquitto stop
mosquitto -c mosquitto.conf