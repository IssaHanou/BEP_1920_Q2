import os
import time

from cc_library.src.sciler.scclib.app import SccLib
from cc_library.src.sciler.scclib.device import Device


class async_device(Device):
    def get_status(self):
        print("status")

    def perform_instruction(self, action):
        print("instruction: " + action)

    def test(self):
        print("test")

    def __init__(self):
        Device.__init__(self)

    def main(self):

        try:
            device = self
            two_up = os.path.abspath(os.path.join(__file__, ".."))
            rel_path = "./async_device_config.json"
            abs_file_path = os.path.join(two_up, rel_path)
            config = open(file=abs_file_path)
            self.scclib = SccLib(config, device)

            self.scclib.__start()

            i = 0
            while True:
                print(i)
                i += 1
                time.sleep(1)

                if i > 100:
                    break

            self.scclib.__stop()

        except KeyboardInterrupt:
            self.scclib.logger.log("program was terminated from keyboard input")
        finally:
            self.scclib.logger.log("cleanly exited async_device program")
            self.scclib.logger.close()


if __name__ == "__main__":
    device = async_device()
    device.main()
