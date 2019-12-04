import datetime
import os

class Logger:
    """
    Logger that prints to console and write to log file:
    A logs folder is created in roots, and txt files are stored there,
    with name format 'log-dd-mm-yyyy--hh-mm-ss.txt'.
    """

    def __init__(self):
        if not os.path.exists("logs"):
            os.mkdir("logs")

        self.filename = "logs/log-" + datetime.datetime.now().strftime("%d-%m-%YT--%H-%M-%S") + ".txt"
        self.file = open(self.filename, "w+")

    def log(self, text):
        """
        Manual logger for the developers.
        """
        print("python log: ", text)
        self.file.write(text)

    def close(self):
        """
        Close file on exit, called on disconnect.
        """
        self.file.close()

