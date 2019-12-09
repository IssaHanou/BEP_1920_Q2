`message`: this is the format of messages sent between all system components.
 The contents types are specific to client computers in this example. 

- `device_id`: the id of the device. 
- `time_sent`: the time at which the message is sent in the format "dd-mm-yyyy hh:mm:ss".
- `type`: the type of the message, this can be:
    - `status`
    - `confirmation`
    - `connection`
    - `instruction`
- `contents`:
    - If type is "status", then the message contents is:
        - `input`: specifies current values for every input_component (in the form "componentId": "value") or an empty string if there is no input 
        - `output`: specifies current values for every output_component (in the form "componentId": "value") or an empty string if there is no output 
    - If type is "confirmation", then the message contents is:
        - `completed`: "true" or "false" depending on success
        - `instructed`: the instruction message that the client computer acted on regarding this confirmation
    - If type is `connection`, then the message contents is:
        - `connection`: "true" or "false" depending on connection status 
    - If type is "instruction", then the message contents is:
        - `instruction`: this is the type of instruction, which can be one of the following:
            `output`, `test`, `start`, `reset`, `stop`. With output, the next property will be:
        - ` output`: defines values for each output_component (in the form "componentId": "value")``
