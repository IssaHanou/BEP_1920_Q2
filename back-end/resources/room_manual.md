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
- `input`: defines type of values to be expected as input. The keys are component ids and values are types of input (in string format).  
    Possible types are: "string", "boolean", "numeric", "array", or a custom name. 
- `output`: defines type of values to be expected as output. The keys are component ids and the value is a map with a `type` property 
    and an optional `instruction`, which defines a map with custom instruction for the device. 
    
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
- `conditions`: this is either a logical operator (i) defined by `operator` (either `AND` or `OR`) and `list` which is a list of conditions or other logical operators **or** this is a condition (ii) defined by `type`, `type_id` and `constraint`
    
    1. Logical operator
        - `operator`: this can be `AND` or `OR`
        - `list`: this is an array of conditions / logical operators
    2. Condition
        - `type`: this can be `rule`, `timer` or `device`.
        - `type_id`: this will be the id of a timer, rule or device, depending on the type.
        - `constraints`: this is either a logical operator (i) defined by `operator` (either `AND` or `OR`) and `list` which is a list of conditions or other logical operators **or** this is a constraint (ii) defined by `comp`, `value` and `component_id`      
        
            1. Logical operator
                - `operator`: this can `AND` or `OR`
                - `list`: this is an array of constraints / logical operators
            2.
                - `comparison`: this is the type of comparison and can be `eq`, `lt`, `gt`, `contains` , `lte`, `gte`. `eq` will work on all types, `lt`, `gt`, `lte`, `gte` only on numeric, and `contains` only on arrays
                - `value`: this is the value on which the comparison is made. This should be in the same type as specified in the input of the device. 
                If it has custom input, then enter value in preferred type and deal with it on the client.
                In case of "timer" type, it should be in the format "hh:mm:ss"
                - `component_id`: in the case of "device" type, this is the id of the component it triggers.
                In case of "timer" type, this is non-existent. 
- `actions`: this is an array of actions:
        
    - `type`: this can be `device` or `timer`
    - `type_id`: the id of device or timer, depending on type respectively
    - `message`: this defines the `output` message sent. The output defines the type of values to be expected as output. 
        In case of device type, the keys are component ids and the value is a type. An additional "instruction" property may be defined.
        In the case of timer, the message should have `instruction` specified as `stop`, `start`, `add`, `subtract` or `set`, in all cases (except `stop` and `start`), a `value` should also be passed. 
