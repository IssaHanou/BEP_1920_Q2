import os
import time
import logging

from sciler.device import Device

max_msg_per_sec = 10


def get_current_time_in_ms():
    return int(round(time.time() * 1000))

class Pinger(Device):
    def __init__(self):
        two_up = os.path.abspath(os.path.join(__file__, ".."))
        rel_path = "pinger_config.json"
        abs_file_path = os.path.join(two_up, rel_path)
        abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
        config = open(file=abs_file_path)
        super().__init__(config)
        self.reset()

    def get_status(self):
        return {
            "#ping": self.ping,
            "#pong": self.pong,
            "pinging": self.pinging,
            "avg": str(self.avg) + "ms",
            "min": str(self.min) + "ms",
            "max": str(self.max) + "ms"
        }

    def perform_instruction(self, action):
        instruction = action.get("instruction")
        if instruction == "pong":
            self.__pong()
        if instruction == "pinging":
            self.pinging = action.get("value")
        else:
            return False
        return True

    def __ping(self):
        self.ping += 1
        self.timeSent = get_current_time_in_ms()
        self.status_changed()

    def __pong(self):
        self.pong += 1
        latency = get_current_time_in_ms() - self.timeSent
        self.accumulativeLatency += latency
        self.latencies.append(latency)
        if self.pong == 1:  # first ping
            self.min = latency
            self.max = latency
            self.avg = latency
        else:  # regular ping
            if latency < self.min:
                self.min = latency
            if latency > self.max:
                self.max = latency
            self.avg = int(round(self.accumulativeLatency / self.pong))

    def test(self):
        print("test")

    def reset(self):
        self.pinging = False
        self.ping = 0
        self.pong = 0
        self.avg = 0
        self.min = 0
        self.max = 0
        self.timeSent = get_current_time_in_ms()
        self.accumulativeLatency = 0
        self.latencies = []
        self.status_changed()

    def __loop(self):
        min_time_between = 1000 / max_msg_per_sec  # min time between messages in milliseconds
        while True:
            time_between = get_current_time_in_ms() - self.timeSent
            if self.pinging and self.ping == self.pong and min_time_between <= time_between:
                self.__ping()

    def __end(self):
        print(self.latencies)

    def main(self):
        logging.basicConfig(
            level=logging.ERROR,
        )
        self.start(loop=self.__loop, stop=self.__end)


if __name__ == "__main__":
    device = Pinger()
    device.main()
