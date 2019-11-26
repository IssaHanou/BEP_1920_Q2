## Configuration file

This manual will help you write a configuration file for an escape room. 
The file should be written in JSON and all the tags necessary are explained here.
An example can be seen in `example.config.json`. The same format 
  
There are three main components to the file:

- `general`
- `devices`
- `puzzles` 
- `general_events`
- `rules` which are defined for puzzles

### General
This is the general information of the escape room. It includes the following tags: 

- `name`: this is the name of the escape room, this is a string, e.g. "Escape room X". This can be displayed in the front-end, so should be readable and in Dutch. 
- `duration`: this is the duration of the escape room, which should be a string in the format "hh:mm:ss".
- `host`: this is the IP address of the broker through which clients and back-end connect, formatted as a string.
- `port`: this is the port on which the broker runs, formatted as integer. 

### Devices
This will be a list of all devices in the room. Each device is defined as a JSON object with the following properties:

- `id`: this is the id of a device. Write it in camelCase, e.g. "controlBoard".
- `description`: this is optional and can contain more information about the device. This can be displayed in the front-end, so should be readable and in Dutch. 
- `input`: defines type of values to be expected as input as a map. There can be one key `value`, or the keys can be component ids. 
    The value is a map with the `type` property. This is defined as a string and can "string", "boolean", "array", "integer" or a custom name. 
- `output`: defines type of values to be expected as output as a map. There can be one key `value`, or the keys can be component ids. 
    The value is a map with the `type` property. This is defined as a string and can "string", "boolean", "array", "integer" or a custom name.
    It can also carry the `instruction` property which defines a map with custom instruction for the device. 
    
### Puzzles
Puzzles is an array of puzzle objects, which have a 

- `name`: name of puzzle. This can be displayed in the front-end, so should be readable and in Dutch. 
- `rules`: array of rule objects (see below)
- `hints`: array of hints (strings), specific to each puzzle. 
These can be displayed in the front-end, so should be readable and in Dutch. 


### General Events
General events have the following properties:

- `name`: name of event, for example "start"
- `rules`: array of rule objects (see below)

### Rules
Rules are defined by:

- `id`: this is the id of a rule. Write it in camelCase, e.g. "solvingControlBoard".
- `description`: this is optional and can contain more information about the rule. 
This can be displayed in the front-end, so should be readable and in Dutch.
- `limit`: this sets the number of times this rule can be triggered. 
- `conditions`: this is an array of conditions. By putting several constraints in an array within the conditions array, they will be treated as OR conditions. 

    - `type`: this can `rule`, `time` or `device`.
    - `id`: this will be the id of a timer, rule or device, depending on the type.
    - `constraints`: this is an array of constraints. By putting several constraints in an array within the constrains array, they will be treated as OR constraints. 
        
        - `comparison`: this can be "eq", "lt", "gt", "cont", "lte", "gte" 
        - `value`: this is the value on which the comparison is made. 
        - `component_id`: in the case of "device" type, where the device has a non-empty input_components list, this is the id of the component it triggers.
- `actions`: this is an array of actions:
        
    - `type`: this can be `device` or `timer`
    - `id`: the id of device or timer, depending on type respectively
    - `message`: this defines the output message sent. In case of device this can either contain a componentIds with their updated values or a general value (array) to the device. 
     In the case of timer, the message should have `instruction` specified as `stop` or `subtract`, in the latter case, a `value` should also be passed. 
