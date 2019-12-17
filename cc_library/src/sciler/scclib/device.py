from abc import ABC, abstractmethod
from cc_library.src.sciler.scclib.app import SccLib


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
        """
        self.scclib = SccLib(config, self)

    def status_changed(self):
        self.scclib.status_changed()

    def start(self, loop=None, stop=None):
        self.scclib.start(loop, stop)

    @abstractmethod
    def main(self):
        """
        This method should be overridden in the subclass.
        It should initialize the SccLib class with config file and device class.
        It should also add event listeners to GPIO for all input components.
        """
        self.start()
