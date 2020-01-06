`message`: this is the format of messages sent between all system components.
### Base of message
- `device_id`: the id of the device of which de message is send. 
- `time_sent`: the time at which the message is sent in the format 
"dd-mm-yyyy hh:mm:ss".
- `type`: the type of the message
- `contents`: the actual contents the message want to pass on

`device_id` and `time_sent` is defined the same for each message, `type` and `contents`
are specific defined depending on the sender and receiver.

### Back-end to Client Computers
- `type`: the type of the message, this can be:
    - `instruction`
- `contents`:
    - If type is `instruction`, then the message contents is __list__ of instructions 
    that have:
        - `instruction`: one of the instructions specified for this device or 
        component or one of following instructions: `test`, `status update`, `reset`
        - `value`: this is the value (argument) for the instruction (optional)
        - `instructed_by`: this is the id of the device which originally send this instruction (usually front-end)
        - `component_id`: this will be the id of a component in a timer or device, 
                (optional)
   
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
        
### Front-end to Back-end
- `type`: the type of the message, this can be:
    - `instruction`
- `contents`:
    - If type is `instruction`, then the then the message contents have
        - `instruction`: one of following instructions: `test all`, `send status`, `hint`, `start`, `stop`, `reset all`.
     
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
    - If type is `status`, then the then the message contents have
        - `id` of device
        - `status` has a map of `component_id` keys and `status` values
        - `connection` boolean
    - If type is `time`, then the then the message contents have
        - `id` of timer
        - `duration` has a number of the duration left in milliseconds
        - `state` sting of the timer state
    - If type is `instruction`, then the then the message contents have
        - `instruction` with value `reset` or `status update`
    - If type is `name`, the contents contains a `name` parameter carrying the name of the escape room
    

        