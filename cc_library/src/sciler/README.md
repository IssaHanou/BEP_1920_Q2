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
- test the device script, by altering `device_test.py`, make sure the class is imported, change the device class in line 13 and the config file name in line 17



## Set-up testing device (full cc_library)
- create custom device script, which should inherit from `Device`, add it to `cc_library/src/scripts`
- write configuration for the the device, according to `device_manual.md`, in the same file as the script
- `ssh pi@<ip-address>`
- `python3 cc_library/src/scripts/<custom-device>.py` to start device
- start broker for device to connect with. run `mosquitto -c <conf>`

## Set-up library
To run this library on a client computer:

- `pip install sciler`
- create custom device script, which should inherit from `Device`, add it to `/scclib`
- write configuration for the the device, according to `device_manual.md`
- `ssh pi@<ip-address>`
- `python3 sclier/scclib/<custom-device>.py` to start device
- start broker for device to connect with. run `mosquitto -c <conf>`

## License
This library is licensed with GNU GPL v3, see `LICENSE.md`.