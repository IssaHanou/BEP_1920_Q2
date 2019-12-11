# BEP_1920_Q2 - broker


### Set-up
When using Mosquitto as a broker, mosquitto.conf should be used to setup a broker with an extra listener for websockets.
 
The bind-address should be changed to local-ip when communication outside a one machine develepment setup is required. 

This command can only be used in the installation folder of Mosquitto. 

The mosquitto.conf part of the command should be replaced by a path to this file.


```
mosquitto -c <mosquitto.conf>
```

When there is a mosquitto broker already running, this will not work. Then, first run `net stop mosquitto`.

### Running
```
mosquitto -t <topic>
mosquitto -h <ip.ad.dr.ess>
```

### Topics

All messages to the front-end should be sent on topic `front-end`

All messages to the back-end should be sent on topic `back-end`

All messages to all client computers should be sent on topic `client-computers`

All messages to single client computers should be sent to topic `<devicename>`

All messages to labeled client computers should be sent to topic `<label>` (like `hint`)
