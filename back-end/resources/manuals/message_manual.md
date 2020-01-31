`message`: this is the format of messages sent between all system components.
### Base of message
- `device_id`: the id of the device of which de message is sent. 
- `time_sent`: the time at which the message is sent in the format 
"dd-mm-yyyy hh:mm:ss".
- `type`: the type of the message
- `contents`: the actual contents the message want to pass on

`device_id` and `time_sent` is defined the same for each message, `type` and `contents`
are specifically defined depending on the sender and receiver.

### Back-end to Client Computers
- `type`: the type of the message, this can be:
    - `instruction`
- `contents`:
    - If type is `instruction`, then the message contents is __list__ of instructions 
    that have:
        - `instruction`: one of the instructions specified for this device or 
        component or one of the following instructions: `test`, `status update`, `reset`
        - `value`: this is the value (argument) for the instruction (optional)
        - `instructed_by`: this is the id of the device which originally send this instruction (usually front-end)
        - `component_id`: this will be the id of a component in a timer or device, 
                (optional)
     
##### example
    { 
    "device_id": "back-end",
    "time_sent": "17-1-2019 16:20:21",
    "type": "instruction",
    "contents": [
            {
            "instruction":"turnOnOff",
            "value": true,
            "component_id": "led1" 
            },
            {
            "instruction":"turnOnOff",
            "value": false,
            "component_id": "led2"
            }
        ]
    }
   
### Client Computers to Back-end
- `type`: the type of the message, this can be:
    - `status`
    - `confirmation`
    - `connection`
- `contents`:
    - If type is `status`, then the message contents is map of status's 
    that have a `key` with a `value`. where `key` is a component_id, and `value` it's status in 
    a format defined in the configurations of the device. e.g. `{redSwitch: true, redSlider: 40, redLed: "aan"}`
    - If type is `confirmation`,  then the message contents has te following:
        - `completed` is a boolean.
        - `instructed` is the original instruction message for the device, including the `instructed_by` tag
    - If type is `connection`, then the message contents has te following:
        - `connection` is a boolean defining the connection status of the device.

##### Example   
    { 
    "device_id": "controlBoard",
    "time_sent": "17-1-2019 16:20:20",
    "type": "status",
    "contents": {
        "redSwitch": true 
        "blueSwitch": false
        }
    }

### Front-end to Back-end
- `type`: the type of the message, this can be:
    - `instruction`
- `contents`:
    - If type is `instruction`, then the message contents have
        - `instruction`: one of following instructions: 
            - `send setup`, `reset all`, `test all`: no additional parameters
            - `test device`, has additional parameter `device` to specify which device to test
            - `finish rule`, has additional parameter `rule` to specify which rule to finish
            - `hint`: has extra parameters `value` and `topic` to specify the hint to send and on which topic

##### Example
    { 
    "device_id": "front-end",
    "time_sent": "17-1-2019 16:20:20",
    "type": "instruction",
    "contents": [{
            "instruction":"finish rule",
            "rule": "playlist puzzel 1",
            }]
    }     

### Back-end to Front-end
- `type`: the type of the message, this can be:
    - `confirmation`
    - `status`
    - `time`
    - `instruction`
- `contents`:
    - If type is `confirmation`, then the then the message contents have
        - `completed`
        - `instructed`, which contains the original instruction with in the `contents`:
            - `instruction`
            - `instructed_by`
    - If type is `instruction`, then the then the message contents have
             - `instruction` with value `reset` or `status update` or `test` or `setState`
    - If type is `status`, then the message contents has
        - `id` of device
        - `status` has a map of `component_id` keys and `status` values
        - `connection` boolean
    - If type is `event status`, then the message contents has objects with
        - `id` of rule
        - `status` of rule describes whether the rule is finished or not
    - If type is `front-end status`, then the message contents has a list of objects with:
        - `id` of button
        - `disabled` status of button
    - If type is `time`, then the message contents have
        - `id` of timer
        - `duration` has a number of the duration left in milliseconds
        - `state` sting of the timer state
    - If type is `setup`, the contents:
        - a `name` parameter carrying the name of the escape room 
        - a `hints` parameter carrying a map with the name of the puzzle as key and a list of hints as value
        - an `events` parameter carrying a list of maps with the parameters `id`, `description`, `status`,
         `puzzle` (whether it belongs to puzzle or not), and `eventName` (name of general event or puzzle it belongs to) 
        - a `cameras` parameter carrying a list of camera objects with a `name` and `link` tag
        - a `buttons` parameter carrying a list of button maps with `id` and `disabled` parameters 
        for buttons that should be added to the front-end
        - a `devices` parameter carrying a list of device maps with `id`, `description` and `labels`

##### example
    { 
    "device_id": "back-end",
    "time_sent": "17-1-2019 16:20:20",
    "type": "status",
    "contents": [
            { 
            "id": "controlBoard" 
            "connection":true
            "status": {
                "redSwitch": true
                "blueSwitch": false
                }
            }
        ]
    }
    { 
    "device_id":"back-end",
    "time_sent":"17-01-2020 15:33:28",
    "type":"event status",
    "contents":[
        {"id":"rood","status":false},
        {"id":"oranje","status":false},
        {"id":"add 5 min timer","status":false},
        {"id":"time up","status":false},
        {"id":"Puzzles done","status":false},
        {"id":"stop timer","status":false},
        {"id":"groen","status":false}
        ]
    }
