# scclib

## Client computer library for S.C.I.L.E.R. system
This mainly consists of:

- `scclib/scclib.py` the main file handling all communication with the rest of the system.
- `scclib/device.py` the superclass from which all client computer handlers should inherit. It defines the three main methods that need to be implemented.
- `scclib/device_test.py` the test file to test an implemented device script with.
- `device_manual.md` the manual for writing configuration files for devices
- `LICENSE.md` the license with which this library complies

## Set-up Device
- run `pip install -i https://test.pypi.org/simple/ scclib`
- create a custom device script, with a class inheriting from the Device superclass, whose main method is called in the script
- write configuration for the the device, according to device_manual.md, in the same folder as the script
- start broker for device to connect with
- run `python3 <custom-device>.py`

## License
This library is licensed with GNU GPL v3, see `LICENSE.md`.


