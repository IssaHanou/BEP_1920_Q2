### Set-up
Make sure you have Python 3 running. Configure the project interpreter as a Python 3 interpreter (Settings -> Project Interpreter).

To make sure all dependencies are loaded, run `pip install -r cc_library/requirements.txt`.

### Checks
- To check the codestyle, run `flake8 cc_library`.
- To check unittests, run `python -m unittest discover`.
- To update formatting, run `black cc_library`.