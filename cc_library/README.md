### Set-up
Make sure you have Python 3 running. Configure the project interpreter as a Python 3 interpreter (Settings -> Project Interpreter).

To make sure all dependencies are loaded, run `pip install -r cc_library/requirements.txt`.

### Checks
- To check unittests, run `python -m unittest discover`.
- To update formatting, run `black cc_library`.
- To check the codestyle, run `flake8 cc_library`.

### Raspberry Pi Set up
- Code specific methods of the device into `<script for device>.py`
- Download `cc_library` onto the Pi
- Add `/home/pi/cc_library/src` to PYTHONPATH
- run `pip3 install -r cc-library/requirements.txt`
- run `cd cc_library`
- run `python3 src/scripts/<script for device>.py`
