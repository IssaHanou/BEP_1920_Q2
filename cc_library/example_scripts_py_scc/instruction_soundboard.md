# Soundboard specifics

#### Messages to send:

status updates:
```json
{ 
    "device_id": "SoundBoard",
    "time_sent": "17-1-2019 16:20:20",
    "type": "status",
    "contents": {
        "knop1": "green",
        "knop2": "off",
        "knop3": "off",
        "knop4": "orange",
        "knop5": "purple",
        "knop6": "off",
        "knop7": "red",
        "knop8": "blue"
        }
}
```
connection message:
```json
{ 
    "device_id": "SoundBoard",
    "time_sent": "17-1-2019 16:20:20",
    "type": "connection",
    "contents": {
        "connection": true 
        }
}
```
see other manual for more

#### Messages to expect:
When the system has to start
```json
 { 
  "device_id": "back-end",
  "time_sent": "17-1-2019 16:19:70",
  "type": "instruction",
  "contents": [
      {
      "instruction":"start",
      "value": true, 
      "component_id": "ready"
      }
  ]
}
```
When the system is done because the puzzle is solved (needs to enable party ending)
```json
{
  "device_id": "back-end",
  "time_sent": "17-1-2019 16:19:70",
  "type": "instruction",
  "contents": [
      {
      "instruction":"correct",
      "value": true, 
      "component_id": "ready"
      }
  ]
}
```
When the system has to give away a hint 
```json
{
  "device_id": "back-end",
  "time_sent": "17-1-2019 16:19:70",
  "type": "instruction",
  "contents": [
      {
      "instruction":"voorzeggen",
      "value": "button5", 
      "component_id": "voorzeggen"
      }
  ]
}
```
    

##### room_config back-end device:
```json
{
  "description": "soundboard",
  "id": "SoundBoard",
  "input": {
    "knop1 kleur": "string",
    "knop2 kleur": "string",
    "knop3 kleur": "string",
    "knop4 kleur": "string",
    "knop5 kleur": "string",
    "knop6 kleur": "string",
    "knop7 kleur": "string",
    "knop8 kleur": "string"
  },
  "output": {
    "ready": {
      "instructions": {
        "start": "boolean",
        "correct": "boolean"
      },
      "type": "string"
    },
    "voorzeggen": {
      "instructions" : {
        "voorzeggen": "string"
      }
    }
  }
}
```