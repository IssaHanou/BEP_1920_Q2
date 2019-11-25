# BEP_1920_Q2 - broker

When using Mosquitto as a broker, mosquitto.conf should be used to setup a broker with an extra listener for websockets. 
Bind-address should be changed to local-ip when communication outside a one machine develepment setup is required. This command can only be used in the installation folder of Mosquitto. The mosquitto.conf part of the command should be replaced by a path to this file.

```
mosquitto -c mosquitto.conf
```