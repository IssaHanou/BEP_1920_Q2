`message`: this is the format of messages sent between all system components.
 The contents types are specific to client computers in this example. 

- `device_id`: the id of the device. 
- `time_sent`: the time at which the message is sent in the format "dd-mm-yyyy hh:mm:ss".
- `type`: the type of the message, this can be:
    - `status`
    - `confirmation`
    - `connection`
    - `instruction`
    - `user_instruction`
- `contents`:
    - If type is `status`, then the message contents is:
        - `input`: specifies current values for every input_component (in the form "componentId": "value") or an empty string if there is no input 
        - `output`: specifies current values for every output_component (in the form "componentId": "value") or an empty string if there is no output 
    - If type is `confirmation`, then the message contents is:
        - `completed`: "true" or "false" depending on success
        - `instructed`: the instruction message that the client computer acted on regarding this confirmation
    - If type is `connection`, then the message contents is:
        - `connection`: "true" or "false" depending on connection status 
    - If type is `instruction`, then the message contents is list of `componentInstructions` that have:
        - `component_id`: this will be the id of a component in a timer or device
        - `instruction`: one of the instructions specified for this device and component
        - `value`: this is the value for the instruction of the type specified for this device and component
    - If type is `user_instruction`, then the message contents is:
        - `instruction`: this is the type of instruction, which can be one of the following: `test`, `start`, `reset`, `stop`. 
        
