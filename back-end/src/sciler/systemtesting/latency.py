import json
from scclib.device import Device


class Latency(Device):
    def __init__(self):
        config = open(file="latency_config.json")
        super().__init__(config)
        self.ping = False

    def get_status(self):
        return {"ping": self.ping}

    def perform_instruction(self, action):
        instruction = action.get("instruction")
        if instruction == "ping":
            print(action.get("value"))
            self.ping = action.get("value")
            self.status_changed()
        else:
            return False, action
        return True, None

    def test(self):
        print("test")

    def reset(self):
        self.ping = False

    def main(self):
        self.start()


if __name__ == "__main__":
    device = Latency()
    device.main()
