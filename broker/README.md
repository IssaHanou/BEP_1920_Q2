# BEP_1920_Q2 - broker


### Set-up
- run sudo apt-get install mosquitto -y
- run sudo apt-get install mosquitto-client

When using Mosquitto as a broker, mosquitto.conf should be used to setup a broker with an extra listener for websockets.
 
The bind-address should be changed to local-ip when communication outside a one machine development setup is required. 

> This command can only be used in the installation folder of Mosquitto. 

The mosquitto.conf part of the command should be replaced by a path to this file.



```
mosquitto -c <mosquitto.conf>
```

When there is a mosquitto broker already running, this will not work. 
Then, first run `net stop mosquitto` or `sudo service mosquitto stop`.

#### Example mosquitto.conf

```
# Config file for mosquitto
# ============================================
# Default listener
# ============================================
# bind_address ip-address/host name
bind_address 192.168.178.82
# Port to use for the default listener.
port 1883
# Choose the protocol to use when listening.
# This can be either mqtt or websockets
protocol mqtt 
# ============================================
# Extra listeners
# ============================================
# listener port-number [ip address/host name]
listener 8083
# Choose the protocol to use when listening.
# This can be either mqtt or websockets.
protocol websockets

```
### To run on boot
- move mosquitto.conf file to etc/mosquitto/conf.d/
- use tool like supervisord:
- run sudo apt-get install supervisor
- run sudo nano /etc/supervisor/conf.d/broker.conf and save:
```
[program:broker]
command=/bin/sh -c "sudo service mosquitto stop && -c etc/mosquitto/conf.d/mosquitto.conf"
user=pi
group_name=pi
stdout_logfile=/home/pi/sciler_logs/logs_broker.txt
redirect_stderr=true
autostart=true
autorestart=true
```

### Topics

All messages to the front-end should be sent on topic `front-end`

All messages to the back-end should be sent on topic `back-end`

All messages to all client computers should be sent on topic `client-computers`

All messages to single client computers should be sent to topic `<devicename>`

All messages to client computers labeled hint should be sent to topic `hint` 


     
    