# Developing with sciler 

## Client computer library for S.C.I.L.E.R. system
This is the library to create devices to work together with the SCILER system 

### Installation
- pip install with ```pip install sciler```

### Using this library
- import lib with `from sciler.device import Device`
##### create a class that extends `Device`
- in order to do this:
    - implement `getStatus()` which should return a dictionary of the current status
    - implement `performInstruction(action)` which should return a boolean of whether the instruction can be performed, where action has:
        - `instruction`: string with the name of the instruction
        - `value`: any type with a value specific for this instruction
        - `component_id`: string with the name of the component for which the instruction is meant (can be undefined) 
    - implement `test()` which returns nothing, this method should do something visible so the operator can test this device works correctly
    - implement `reset()` which returns nothing, this method should make the device return to its starting state so that the escape room can be started again
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
        - it should also add event listeners to GPIO for all input components.
    - implement `main()` which should call `start(loop, stop)` with an optional event loop and ending function.  
##### Now in your class which implements `Device` 
- you can call:
    - `log(level, message)` which logs using the logger provided in `Device` where level one of the following strings: 'debug', 'info', 'warn', 'error', 'fatal' and message custom string containing more information
    -  `statusChanged()` which can be called to signal to `Device` that the status is changed, this will send a new status to SCILER
##### To now start the system
 - initialize the device in your main program and call `device.main() `

### Example
```python
import os

from sciler.device import Device


class Display(Device):
    def __init__(self):
        two_up = os.path.abspath(os.path.join(__file__, ".."))
        rel_path = "display_config.json"
        abs_file_path = os.path.join(two_up, rel_path)
        abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
        config = open(file=abs_file_path)
        super().__init__(config)
        self.hint = ""

    def get_status(self):
        return {"hint": self.hint}

    def perform_instruction(self, action):
        instruction = action.get("instruction")
        if instruction == "hint":
            self.show_hint(action)
        else:
            return False, action
        return True, None

    def test(self):
        self.hint = "test"
        print(self.hint)
        self.status_changed()

    def show_hint(self, data):
        self.hint = data.get("value")
        print(self.hint)
        self.status_changed()

    def reset(self):
        self.hint = ""
        self.status_changed()

    def main(self):
        self.start()


if __name__ == "__main__":
    device = Display()
    device.main()
```
where `display_config.json` is
```json
{
  "id": "display",
  "description": "Display can print hints",
  "host": "192.168.178.82",
  "labels": ["hint"],
  "port": 1883
}
```
example of `main()` with loop and stop:
```python
    def loop(self):
        previous = self.get_sliders_analog_reading()
        while True:
            positions = self.get_sliders_analog_reading()
            if previous != positions:
                self.status_changed()
                previous = positions

    def main(self):
        self.setup_events()
        self.start(loop=self.loop, stop=GPIO.cleanup)
```

## License
This library is licensed with GNU GPL v3, see `LICENSE.md`.
