import logging
from abc import ABC, abstractmethod
from sciler.app import Sciler


class Device(ABC):
    """
    Abstract device class from which all custom devices should inherit.
    Defines main methods needed for communication to S.C.I.L.E.R.

    A device script should have a class inheriting this class.
    Its main method should be called at the bottom.
    """

    @abstractmethod
    def get_status(self):
        """
        Returns status of all custom components, in a dictionary
        """

    @abstractmethod
    def perform_instruction(self, action):
        """
        Defines how instructions are handled,
        for all instructions defined in output of device in config.
        :param action: a dictionary containing an instruction tag, component_id tag and value tag.
        :return boolean: True if instruction was valid and False if illegal instruction
        was sent or error occurred such that instruction could not be performed.
        Returns tuple, with boolean and None if True and the failed action if false.
        """

    @abstractmethod
    def test(self):
        """
        Defines test sequence for device.
        """

    @abstractmethod
    def reset(self):
        """
        Defines a reset sequence for device.
        """

    def __init__(self, config):
        """
        The init of the subclass should call this method and also initialize all class attributes.
        :param config: the config directory
        """
        self.scclib = Sciler(config, self)

    def status_changed(self):
        """
        Passes the request to send status too the back-end to the message sender
        """
        self.scclib.status_changed()

    def start(self, loop=None, stop=None):
        """
        Passes the request to start the message sender
        :param loop: possible event loop
        :param stop: possible end function
        """
        self.scclib.start(loop, stop)

    def log(self, msg, level=logging.INFO):
        """
        Use the logger to log messages
        :param msg: massage to log
        :param level: level of message
        """
        self.scclib.log(level, msg)

    @abstractmethod
    def main(self):
        """
        This method should be overridden in the subclass.
        It should initialize the SccLib class with config file and device class.
        It should also add event listeners to GPIO for all input components.
        """
        self.start()
