## Configuration file

This manual will help you write a configuration file for an escape room. 
The file should be written in JSON and all the tags necessary are explained here.
An example can be seen in `example.config.json`. The same format 
  
There are three main components to the file:

- `general`
- `devices`
- `timers`
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

- `id`: this is the id of a device. Write it in camelCase, e.g. "controlBoard". This id should be unique compared to other device ids and also the rule ids as well as the timer ids.
- `description`: this is optional and can contain more information about the device. This can be displayed in the front-end, so should be readable and in Dutch. 
- `input`: defines type of values to be expected as input. The keys are component ids and values are types of input (in string format).  
    Possible types are: "string", "boolean", "numeric", "array", or a custom name. 
- `output`: defines what this components outputs as their status and what instructions can be performed on this component
    - `type`: defines type of values to be expected as output. Possible types are: "string", "boolean", "numeric", "array", or a custom name. 
    - `instructions`: this is a map of the name of an instruction to the type of argument the instruction takes
    - `label`: this is a list of possible labels this component listens to when an action gets called on a label.
    
### Timers
This will be a list of all the time related actions/conditions. all timers have to be started in a action and be checked in a condition to be used.
- `id`: this will be the id of the timer. Write in camelCase and numbers are fine, e.g. "timerHint1". This id should be unique compared to other timers ids and also the rule ids as well as the device ids.
- `duration`: This will be the duration after which the timer will trigger to true and the conditions containing the timer will be checked to execute actions. The format is XhXmXs, each size optional, e.g. 1h30m30s, 40m30s, 1m30s

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

- `id`: this is the id of a rule. Write it in camelCase, e.g. "solvingControlBoard". This id should be unique compared to other rule ids and also the device ids as well as the timer ids.
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
                - `comparison`: this is the type of comparison and can be `eq`, `lt`, `gt`, `contains` , `lte`, `gte`, `not`. However, only `eq` will work on all types, `lt`, `gt`, `lte`, `gte` only on numeric, and `contains` only on arrays, `not` does not work on booleans.
                - `value`: this is the value on which the comparison is made. In case of `device` type, it should be in the same type as specified in the input of the device. 
                If it has custom input, then enter value in preferred type and deal with it on the client.
                In case of `timer` type, it should be boolean
                In case of `rule` type, it should be numeric since the comparison will be done against the times the rule is executed
                - `component_id`: in the case of "device" type, this is the id of the component it triggers.
                In case of "timer" type, this is non-existent. 
- `actions`: this is an array of actions:
        
    - `type`: this can be `device`, `timer` or `label`
    - `type_id`: the id of device, timer or label, depending on type respectively
    - `message` in case of type `device`: this defines a list of componentInstructions which have:
        - `component_id`: this will be the id of a component in a timer or device
        - `instruction`: one of the instructions specified for this device and component
        - `value`: this is the value for the instruction of the type specified for this device and component
    - `message` in case of type `timer`:   
           - `instruction`: one of the instructions for timer, e.g. `start`, `stop`, `pause`, `done`, `add`, `subtract`
           - `value`: optional, in case of `add` and `subtract` a time should be given in format XhXmXs 
    - `message` in case of type `label`:   
           - `instruction`: one of the instructions specified for the components with this label
           - `value`: this is the value for the instruction of the type specified for this device and component 
    - `delay` in case of type `device` or `label`: This is optional, this is a duration in format XhXmXs, if an action has a delay, the message will publish after this delay.
