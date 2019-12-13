import os
from abc import abstractmethod
from cc_library.src.sciler.scclib.app import SccLib


class Device:
    """
    Abstract device class from which all custom devices should inherit.
    Defines main methods needed for communication to S.C.I.L.E.R.

    A device script should have a class inheriting this class.
    Its main method should be called at the bottom.
    """

    @abstractmethod
    def get_status(self):
        """
        Returns status of all custom components, in json format.
        """

    @abstractmethod
    def perform_instruction(self, contents):
        """
        Defines how instructions are handled,
        for all instructions defined in output of device in config.
        :param contents: contains instruction tag and calls the appropriate functions.
        :return boolean: True if instruction was valid and False if illegal instruction
        was sent or error occurred such that instruction could not be performed.
        Returns tuple, with boolean and None if True and the failed action if false.
        """

    @abstractmethod
    def test(self):
        """
        Defines test sequence for device.
        """

    def __init__(self):
        """
        The init of the subclass should call this method and also initialize all class attributes.
        """
        self.scclib = None

    @abstractmethod
    def main(self):
        """
        This method should be overridden in the subclass.
        It should initialize the SccLib class with config file and device class.
        It should also add event listeners to GPIO for all input components.
        """
        try:
            device = self

            two_up = os.path.abspath(os.path.join(__file__, ".."))
            rel_path = "<enter-filename.json>"
            abs_file_path = os.path.join(two_up, rel_path)
            abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
            config = open(file=abs_file_path)
            self.scclib = SccLib(config, device)

            """Initialize the event listeners here."""

            self.scclib.start()
        except KeyboardInterrupt:
            self.scclib.logger.log("program was terminated from keyboard input")
        finally:
            self.scclib.logger.log("Cleanly exited ControlBoard program")
            self.scclib.logger.close()
