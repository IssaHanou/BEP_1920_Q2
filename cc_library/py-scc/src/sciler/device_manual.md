Each device is defined as a JSON object with the following properties:

- `id`: this is the id of a device. Write it in camelCase, e.g. "controlBoard".
- `description`: this is optional and can contain more information about the device. This can be displayed in the front-end, so should be readable and in Dutch. 
- `host`: the IP address of the host for the broker, formatted as a string.
- `port`:  the port of the host for the broker, formatted as a number.
- `labels`: a list of the labels belonging to the device to receive labeled instructions.
- `input`: defines type of values to be expected as input as a map. There can be one key `value`, or the keys can be component ids. 
    The value is a map with the `type` property. This is defined as a string and can "string", "boolean", "array", "integer" or a custom name. 
- `output`: defines type of values to be expected as output as a map. There can be one key `value`, or the keys can be component ids. 
    The value is a map with the `type` property. This is defined as a string and can "string", "boolean", "array", "integer" or a custom name.
    It can also carry the `instruction` property which defines a map with custom instruction for the device. 

