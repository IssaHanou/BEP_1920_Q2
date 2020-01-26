#!/bin/bash
echo "updating system"
sudo apt-get update -y

echo "installing mosquitto"
sudo apt-get install -y mosquitto
sudo apt-get install -y mosquitto-clients

echo "stopping any mosquitto service and starting own mosquitto"
sudo service mosquitto stop
mosquitto -c "mosquitto.conf" -d
echo "broker running"