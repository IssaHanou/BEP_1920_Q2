# BEP_1920_Q2 - broker


### Set-up
When using Mosquitto as a broker, mosquitto.conf should be used to setup a broker with an extra listener for websockets.
 
The bind-address should be changed to local-ip when communication outside a one machine develepment setup is required. 

This command can only be used in the installation folder of Mosquitto. 

The mosquitto.conf part of the command should be replaced by a path to this file.


```
mosquitto -c <mosquitto.conf>
```

When there is a mosquitto broker already running, this will not work. 
Then, first run `net stop mosquitto` or `sudo service mosquitto stop`.

### Topics

All messages to the front-end should be sent on topic `front-end`

All messages to the back-end should be sent on topic `back-end`

All messages to all client computers should be sent on topic `client-computers`

All messages to single client computers should be sent to topic `<devicename>`

All messages to labeled client computers should be sent to topic `<label>` (like `hint`)

### Client devices without library

To set up communication with the S.C.I.L.E.R. system without the use of a library, you need to consider the following points:

- Connect with a mosquitto broker at the same IP adress as the rest of the system, use the ports defined in your mosquitto config file or add a port if the current ones are not compatible with your device
- The MQTT client should subscribe to the topic `<devicename>` and `"client-computers"` (and `hint` if it is a hint device), where <devicename> is the same name used in the configfile of the back-end
- A MQTT message publisher should publish messages to topic `back-end`, defined in the message_manual.md under chapter __Client Computers to Back-end__
- On start-up a connection message should be send, and set a will for the mqtt client with a disconnection message, here is a python example:
            
        client.will_set("back-end", 
         json.dumps({
         "device_id": <devicename>,
         "time_sent": datetime.now().strftime("%d-%m-%Y %H:%M:%S"),
         "type": "connection",
         "contents": {"connection": False}
         }),
         1)
         
- A MQTT message handler should be able to handle all messages defined in the message_manual.md under chapter __Back-end to Client Computers__. 
Make sure your message handler can process all messages or check if it can process it. after a message, publish a confirmation message back to back-end,
python example:

        message = message.payload.decode("utf-8")
        message = json.loads(message)
        success = self.__check_message(message.get("contents"))
        conf_msg_dict = {
                    "device_id": self.name,
                    "time_sent": datetime.now().strftime("%d-%m-%Y %H:%M:%S"),
                    "type": "confirmation",
                    "contents": {"completed": success, "instructed": message},
                }
        msg = json.dumps(conf_msg_dict)
        client.publish("back-end", msg, 1)
     
    