This will be a list of all devices in the room. Each device is defined as a JSON object with the following properties:

- `id`: this is the id of a device. Write it in camelCase, e.g. "controlBoard".
- `description`: this is optional and can contain more information about the device. This can be displayed in the front-end, so should be readable and in Dutch. 
- `IP_address`: the IP address of the client computer for this device
- `input_components`: this is a list of ids of different components of a device, which can generate input to the system. This can be empty. 
- `output_components`: this is a list of ids of different components of a device, which can provide output of the system. This can be empty. 
- `message`: this is an example message of the status update. It contains two status messages:
    
    - `input`: defines either default values for every input_component (in the form "componentId": "value") or if the input_components list is empty it will define the `value` tag, which can be a string, integer, boolean or array of single values.
    - `output`: definer either values for every output_component (in the form "componentId": "value") or if the output_components list is empty it will define the `value` tag, which can be a string, integer, boolean or array of single values.
- `test`: defines test output sequence in the form of an `output` tag, defined similarly as described above. 
- `label`: *To be added later*
- `interval`: *To be added later*
