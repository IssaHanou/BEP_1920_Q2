- `message`: this is an example message of the status update. Which is a separate JSON object:
    - `message_id`: the id of the message, which should be generated.
    - `device_id`: the id of the device. 
    - `time_sent`: the time at which the message is sent in the format "dd-mm-yyyy hh:mm:ss".
    - `type`: the type of the message, this can be:
        - `status`
        - `confirmation`
        - `instruction`
    - `contents`:
        - If type is "status", then the message contents contains a status: 
            - If type is "component", then this contains the default (reset) settings of the device, for each of the components separately. 
            If the status is numeric, it can be entered as a number. If it is a boolean condition (e.g. on/off), then it can be entered as "true" or "false". Else, enter a string.
            - If type is "sequence", then the status will be an array of a sequence example, e.g. [0, 3, 2, 1].
        - If type is "confirmation", then the message contents is:
            - `completed`: "true" or "false" depending on success
        - If type is "instruction", then the message contents is:
            - `instruction`: this is the type of instruction, which can be one of the following:
            `update`, `test`, `start`, `reset`, `stop`. With update, the `new_status` will be added below instruction 