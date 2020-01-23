```
 ________  ________  ___  ___       _______   ________     
|\   ____\|\   ____\|\  \|\  \     |\  ___ \ |\   __  \    
\ \  \___|\ \  \___|\ \  \ \  \    \ \   __/|\ \  \|\  \   
 \ \_____  \ \  \    \ \  \ \  \    \ \  \_|/_\ \   _  _\  
  \|____|\  \ \  \____\ \  \ \  \____\ \  \_|\ \ \  \\  \| 
    ____\_\  \ \_______\ \__\ \_______\ \_______\ \__\\ _\ 
   |\_________\|_______|\|__|\|_______|\|_______|\|__|\|__|
   \|_________|                                            
```                                                           

### Structure
There are two versions of the S.C.I.L.E.R. Client Computer Library:
   - js-scc: a JavaScript version installable by npm, [click here for the  `README.md`](js-scc/README.md)
   - py-scc: a Python version (will be) installable by pip [click here for the  `README.md`](py_scc/README.md)

### Client devices without library

To set up communication with the S.C.I.L.E.R. system without the use of a library, you need to consider the following points:

- Connect with a mosquitto broker at the same IP adress as the rest of the system, use the ports defined in your mosquitto config file or add a port if the current ones are not compatible with your device
- The MQTT client should subscribe to the topic `<devicename>` and `"client-computers"` (and `hint` if it is a hint device), where <devicename> is the same name used in the configfile of the back-end
- A MQTT message publisher should publish messages to topic `back-end`, defined in the message_manual.md under chapter __Client Computers to Back-end__
- On start-up a connection message should be send, and set a will for the mqtt client with a disconnection message, here is a python example:
```python
client.will_set("back-end", json.dumps({
   "device_id": <devicename>,
   "time_sent": datetime.now().strftime("%d-%m-%Y %H:%M:%S"),
    "type": "connection",
    "contents": {"connection": False} }),
    1)
```
- A MQTT message handler should be able to handle all messages defined in the message_manual.md under chapter __Back-end to Client Computers__. 
Make sure your message handler can process all messages or check if it can process it. after a message, publish a confirmation message back to back-end,
python example:
```python
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
```
       