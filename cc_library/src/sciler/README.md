# scclib

## Client computer library for S.C.I.L.E.R. system
This mainly consists of:

- `scclib/scclib.py` the main file handling all communication with the rest of the system.
- `scclib/device.py` the superclass from which all client computer handlers should inherit. It defines the three main methods that need to be implemented.
- `scclib/device_test.py` the test file to test an implemented device script with.
- `device_manual.md` the manual for writing configuration files for devices
- `LICENSE.md` the license with which this library complies

## Set-up device with package.
- create custom device script, which should inherit from `Device`, add it to `cc_library/src/scripts`
- write configuration for the the device, according to `device_manual.md`, add it to `cc_library/src/scripts`
- test the device script, by altering `device_test.py`, import the class on line 7, change the device class in line 15 and the config file name in line 21



## Set-up testing device (full cc_library)
- create custom device script, which should inherit from `Device`, add it to `cc_library/src/scripts`
- write configuration for the the device, according to `device_manual.md`, in the same file as the script
- start broker for device to connect with. 
- move the cc_library directory onto the Pi
- run `pip3 install paho-mqtt` on Pi
- run `python3 cc_library/src/scripts/<custom-device>.py` on Pi to start device

## Set-up Pi (TODO)


## License
This library is licensed with GNU GPL v3, see `LICENSE.md`.