### Set-up
Make sure you have Python 3 running. Configure the project interpreter as a Python 3 interpreter (Settings -> Project Interpreter).

To make sure all dependencies are loaded, run `pip install -r cc_library/requirements.txt`.

### Checks
- To check unittests, run `python -m unittest discover`.
- To update formatting, run `black cc_library`.
- To check the codestyle, run `flake8 cc_library`.

### Raspberry Pi Set up
- Code specific methods of the device into start_device.py
- Download cc_library/src onto the Pi
- Add /home/pi/src to PYTHONPATH
- pip install all requirements
- run python src/scripts/start_device.py
