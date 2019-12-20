### Set-up
Make sure you have Python 3 running. Configure the project interpreter as a Python 3 interpreter (Settings -> Project Interpreter).

To make sure all dependencies are loaded, run `pip install -r cc_library/requirements.txt`.

### Checks
- To check unittests, run `python -m unittest discover`.
- To update formatting, run `black cc_library`.
- To check the codestyle, run `flake8 cc_library`.

### Raspberry Pi Set up
For setting up a Pi, [a read the `README.md`](src/sciler/README.md) in the sciler package. 