from abc import abstractmethod


class Device:
    """
    Abstract device class from which all custom devices should inherit.
    Defines three main methods needed for communication to S.C.I.L.E.R.
    """

    @abstractmethod
    def get_status(self):
        """
        Returns status of all custom components, in json format.
        """

    @abstractmethod
    def perform_instruction(self, contents):
        """
        Defines how instructions are handled.
        :param contents: contains instruction tag and calls the appropriate functions.
        """

    @abstractmethod
    def test(self):
        """
        Defines test sequence for device.
        """
