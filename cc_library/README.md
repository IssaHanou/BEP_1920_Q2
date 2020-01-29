# Client Computer libraries

## Deployment
There are two versions of the sciler Client Computer Library:
   - a JavaScript version installable by npm, [click here for the  `README.md`](js_scc/README.md)
   - a Python version (will be) installable by pip [click here for the  `README.md`](py_scc/README.md)

## Development without library

To set up communication with the S.C.I.L.E.R. system without the use of a library, you need to consider the following points:

- Connect with a mosquitto broker at the same IP address as the rest of the system, use the ports defined in your mosquitto config file or add a port if the current ones are not compatible with your device
- The MQTT client should subscribe to the topic `<devicename>` and `"client-computers"` (and `hint`/`time` if it is a hint/time device), where <devicename> is the same name used in the configfile of the back-end
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
       
status message: 
```json
{ 
    "device_id": "controlBoard",
    "time_sent": "17-1-2019 16:20:20",
    "type": "status",
    "contents": {
        "redSwitch": true,
        "blueSwitch": false
        }
}
```
connection message: 
```json
{ 
    "device_id": "controlBoard",
    "time_sent": "17-1-2019 16:20:20",
    "type": "connection",
    "contents": {
        "connection": true 
        }
}
```
confirmation message:
```json
{ 
    "device_id": "controlBoard",
    "time_sent": "17-1-2019 16:20:21",
    "type": "confirmation",
    "contents": {
        "completed": true, 
        "instructed": { 
          "device_id": "back-end",
          "time_sent": "17-1-2019 16:19:70",
          "type": "instruction",
           "contents": [
              {
              "instruction":"test",
              "instructed_by": "front-end" 
              }
           ]
        }
    }
    }
```
instruction message:
```json
{ 
          "device_id": "back-end",
          "time_sent": "17-1-2019 16:19:70",
          "type": "instruction",
           "contents": [
              {
              "instruction":"turnOnOFf",
              "value" : true,
              "component_id": "redled"
              }
           ]
        }
```
time message:
```json
{ 
          "device_id": "back-end",
          "time_sent": "17-1-2019 16:19:70",
          "type": "time",
           "contents": [
              {
              "id":"general",
              "duration" : 500000,
              "state": "stateActive"
              }
           ]
        }
```
## Development

### sciler Javascript
For checkstyle prettier is used.
#### Installation prettier
```
    npm install --save-dev --save-exact prettier
    npm install --save-dev tslint-config-prettier
```
#### Run prettier manually
- navigate to `BEP_1920_Q2/front-end/`
- run `prettier --write "**/*.ts"`

#### uploading library to npm
in `package.json` change version, then run:
```
npm login
npm publish
```

### sciler Python
For checkstyle flake8 and black are used.
#### Installation black and flake8
- run `pip install black flake8`
#### Run black and flake8 manually
- navigate to `cc_library\py_scc\`
- run `black sciler`
- run `flake8 sciler`

#### uploading library to pip
in setup.py change version, then run:
```
pip install twine wheel 
python setup.py sdist bdist_wheel
twine upload dist/*
```
