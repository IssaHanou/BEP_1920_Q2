- `message`: this is an example message of the status update. Which is a separate JSON object:
    - `device_id`: the id of the device. 
    - `time_sent`: the time at which the message is sent in the format "dd-mm-yyyy hh:mm:ss".
    - `type`: the type of the message, this can be:
        - `status`
        - `confirmation`
        - `connection`
        - `instruction`
    - `contents`:
        - If type is "status", then the message contents is:
            - `input`: defines either default values for every input_component (in the form "componentId": "value") 
                or if the input_components list is empty it will define the `value` tag, which can be a string, integer, boolean or array of single values.
        - If type is "confirmation", then the message contents is:
            - `completed`: "true" or "false" depending on success
        - If type is `connection`, then the message contents is:
            - `connection`: "true" or "false" depending on connection status 
        - If type is "instruction", then the message contents is:
            - `instruction`: this is the type of instruction, which can be one of the following:
                `output`, `test`, `start`, `reset`, `stop`. With output, the next property will be:
            - ` output`: defines either values for every output_component (in the form "componentId": "value")
                or if the output_components list is empty it will define the `value` tag, which can be a string, integer, boolean or array of single values.
 