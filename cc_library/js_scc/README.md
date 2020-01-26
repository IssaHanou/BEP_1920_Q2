## sciler [(see npm)](https://www.npmjs.com/package/sciler)

### Installation
- npm install with ```npm install sciler```

### Using this library
- create a class that extends `Device`, in order to do this, implement the following methods:

    | method                       | parameters                       |                                                                                                                  | returns                                                                                                                              |
    |------------------------------|----------------------------------|------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------|
    | `getStatus()`                | none                             |                                                                                                                  | a dictionary of current status of the device, this should include all input and output as defined in the room_config for the backend |
    | `performInstruction(action)` | action (a dictionary with keys:) | `instruction` (string with the name of the instruction)                                                          | boolean of  whether the instruction could be performed                                                                               |
    |                              |                                  | `value` (any type with a value specific for this instruction)                                                    |                                                                                                                                      |
    |                              |                                  | `component_id` (string with the name of the component for which the instruction is meant, this can be undefined) |                                                                                                                                      |
    | `test()`                     | none                             |                                                                                                                  | void                                                                                                                                 |
    | `reset()`                    | none                             |                                                                                                                  | void                                                                                                                                 |

- create a constructor which calls the constructor of `Device` with `super(config, logger)` where:
    - config is a dictionary which has keys:
    
   | key      | explanation                                                                |
   |----------|----------------------------------------------------------------------------|
   | `id`     | string with the id of a device. Write it in camelCase, e.g. "controlBoard" |
   | `host`   | string with the IP address of the host for the broker                      |
   | `port`   | int with the port of the host for the broker                               |
   | `labels` | array of strings with all labels this device should also subscribe to      |
   
    - logger is a function(date, level, message) in which an own logger is implemented where
         - `date` is an Date object
         - `level` is one of the following strings: 'debug', 'info', 'warn', 'error', 'fatal'
         - `message` is a custom string containing more information
- Now on the instantiation of your class which implements `Device`, you can call:

    | function              | arguments                                                                         | returns | usecase                                                                          |
    |-----------------------|-----------------------------------------------------------------------------------|---------|----------------------------------------------------------------------------------|
    | `start(onStart)`      | onStart (a function that gets called once the device is connected or reconnected) | void    | call this function in order to connect to sciler                                 |
    | `log(level, message)` | level (one of the following strings: 'debug', 'info', 'warn', 'error', 'fatal')   | void    | this method can be used to log in the same logger as the library does            |
    |                       | message (custom string containing more information)                               |         |                                                                                  |
    | `statusChanged()`     | none                                                                              | void    | call this function to signal that you updated a status so sciler can be notified |

- import lib with ```const Device = require("sciler");```
- in case of:
    - angular: add `sciler` to dependencies in `package.json`
    - browser javascript: (example nodejs serving web page with javascript which includes this library), [Browserify](http://browserify.org/) your javascript which includes this library.

### Example
Javascript file:
```javascript
$(document).ready(function() {
  const Device = require("sciler");
  let display;

  class Display extends Device {
    constructor(config) {
      super(config, timedLogger);
      this.hint = "";
      this.button = false;
    }

    // required method for extending Device
    getStatus() {
      return {
        button: this.button,
        hint: this.hint
      };
    }

    // required method for extending Device
    performInstruction(action) {
      switch (action.instruction) {
        case "hint": {
          this.hint = action.value;
          displayText(this.hint);
          this.statusChanged();
          break;
        }
        default: {
          return false;
        }
      }
      return true;
    }

    // required method for extending Device
    test() {
      this.hint = "test";
      displayText(this.hint);
    }

    // required method for extending Device
    reset() {
      this.hint = "";
      this.button = false;
      displayText(this.hint);
      this.statusChanged();
    }
  }

  // custom logger used in constructor
  function timedLogger(date, level, message) {
    const formatDate = function(date) {
      return (
        date.getDate() +
        "-" +
        date.getMonth() +
        1 +
        "-" +
        date.getFullYear() +
        " " +
        date.getHours() +
        ":" +
        date.getMinutes() +
        ":" +
        date.getSeconds()
      );
    };
    console.log(
      "time=" + formatDate(date) + " level=" + level + " msg=" + message
    ); // call own logger);
  }

  // edit the DOM using JQuery to display text
  function displayText(text) {
    $("#hint").text(text);
  }

  // when the button is click update status and notify sciler
  $("#button").on("click", function() {
    display.button = true;
    display.statusChanged();
  });

  // get config file from server
  $.get("/display_config.json", function(config) {
    display = new Display(JSON.parse(config)); // create new Display object
    // connect
    display.start(() => {
      console.log("connected"); // when connected, do something
    });
  });
});
```
Where `display_config.json` is:
```json
{
  "id": "display-node",
  "host": "192.168.178.49",
  "labels": ["hint"],
  "port": 8083
}
```
And where `index.html` is:
```haml
<!DOCTYPE html>
<html lang="en" xmlns="http://www.w3.org/1999/html">
<head>
    <meta charset="UTF-8">
    <title>Display</title>
    <link rel="icon" type="image/x-icon" href="raccoon.ico">
    <script type="text/javascript" src="http://code.jquery.com/jquery-1.7.1.min.js"></script>
    <script type="text/javascript" src="bundle.js"></script>
</head>
<body>
hint: <p id="hint"></p>
<button id="button">button</button>
</body>
</html>
```
And where the `room_config.json` looks something like: 
```json
{
  "general": { },
  "cameras": [ ],
  "general_events": [ ],
  "button_events": [ ],
  "puzzles": [ ],
  "timers": [ ],
  "devices": [
    { },
    {
      "id": "display-node",
      "description": "displays messages",
      "input": {
        "button": "boolean"
      },
      "output": {
        "display": {
          "type": "string",
          "instructions": {
            "hint": "string"
          }
        }
      }
    },
    { }
  ]
}
```

### Tip reading in config.json in Angular:
When reading in a json file
```
import * as data from './data.json';
```

make sure `tsconfig.json` has `resolveJsonModule: true`:
```
{
  "compilerOptions": {
    ...
    "resolveJsonModule": true,
    ...
}
```

then you can use this object in for example `ngOnInit`
```
  ngOnInit(): void {
       this.jsonData = (data  as  any).default
  }
```

### License
GNU GENERAL PUBLIC LICENSE Version 3, [see `LICENSE.md`](LICENSE.md)
