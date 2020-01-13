                                                         
To make sure all dependencies are loaded, run `pip install -r cc_library/requirements.txt`.

### Checks
- To check unittests, run `python -m unittest discover`.
- To update formatting, run `black py-scc`.
- To check the codestyle, run `flake8 py-scc`.

### Raspberry Pi Set up
For setting up a Pi, [a read the `README.md`](src/sciler/README.md) in the sciler package. 


## Set-up Device
- run `pip install -i https://test.pypi.org/simple/ scclib`
- create custom device script, which should inherit from Device
- write configuration for the the device, according to device_manual.md, in the same folder as the script
- start broker for device to connect with
- run `python3 <custom-device>.py`
 
 
## Set-up testing device (full cc_library)
- create custom device script, which should inherit from `Device`, add it to `cc_library/src/scripts`
- write configuration for the the device, according to `device_manual.md`, in the same file as the script
- start broker for device to connect with. 
- move the cc_library directory onto the Pi
- run `pip3 install paho-mqtt` on Pi
- run `python3 cc_library/src/scripts/<custom-device>.py` on Pi to start device

## Set-up device with package.
- create custom device script, which should inherit from `Device`, add it to `cc_library/src/scripts`
- write configuration for the the device, according to `device_manual.md`, add it to `cc_library/src/scripts`
- test the device script, by altering `device_test.py`, import the class on line 7, change the device class in line 15 and the config file name in line 21
