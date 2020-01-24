## js-scc [(see npm)](https://www.npmjs.com/package/js-scc)

### Installation
- npm install with ```npm install js-scc```

### Using this library
- import lib with ```const Device = require("js-scc");```
- create a class that extends `Device`, in order to do this:
    - implement getStatus() which should return a dictionary of the current status
    - implement performInstruction(action) which should return a boolean of whether the instruction can be performed, where action has:
        - `instruction`: string with the name of the instruction
        - `value`: any type with a value specific for this instruction
        - `component_id`: string with the name of the component for which the instruction is meant (can be undefined) 
    - implement test() which returns nothing, this method should do something visible so the operator can test this device works correctly
    - implement reset() which returns nothing, this method should make the device return to its starting state so that the escape room can be started again
    - create a constructor which calls the constructor of `Device` with `super(config, logger)` where:
        - config is a dictionary which has keys:
            - `id`: this is the id of a device. Write it in camelCase, e.g. "controlBoard".
            - `host`: the IP address of the host for the broker, formatted as a string.
            - `port` the port of the host for the broker, formatted as a number.
            - `labels`: these are the labels to which this device should also subscribe, labels is an array of strings, 
        - logger is a function(date, level, message) in which an own logger is implemented where
             - `date` is an Date object
             - `level` is one of the following strings: 'debug', 'info', 'warn', 'error', 'fatal'
             - `message` is a custom string containing more information
- Now in your class which implements `Device`, you can call:
    - start(onStart) which connects the device where onStart is a function that will be called once the device is connected
    - log(level, message) which logs using the logger provided in `Device` where level one of the following strings: 'debug', 'info', 'warn', 'error', 'fatal' and message custom string containing more information
    - statusChanged() which can be called to signal to `Device` that the status is changed, this will send a new status to SCILER
- in case of:
    - angular: add `js-scc` to dependencies in `package.json`
    - browser javascript: (example nodejs serving web page with javascript which includes this library), [Browserify](http://browserify.org/) your javascript which includes this library.

### Example
Javascript file:
```javascript
$(document).ready(function() {
  const Device = require("js-scc");
  let display;

  class Display extends Device {
      constructor(config) {
          super(config, function (date, level, message) {
              const formatDate = function (date) {
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
          });
          this.hint = "";
      }

      getStatus() {
          return {
              "hint": this.hint
          }
      }

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

      test() {
          this.hint = "test";
          displayText(this.hint);
      }

      reset() {
          this.hint = "";
          displayText(this.hint);
          this.statusChanged();
      }
  }

  function displayText(text) {
      $("#hint").text(text);
  }


  // get config file from server
  $.get("/display_config.json", function(config) {
      display = new Display(JSON.parse(config));
      display.start(() => {
          console.log("connected");
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
<p id="hint"></p>
</body>
</html>


```


### Tip reading in config.json in Angular:
When reading in a json file
```
import * as data from './data.json';
```

make sure `tsconfig.json` hase `resolveJsonModule: true`:
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
