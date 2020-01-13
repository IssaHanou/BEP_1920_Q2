import os

from cc_library.src.sciler.scclib.device import Device


class Display(Device):
    def __init__(self):
        two_up = os.path.abspath(os.path.join(__file__, ".."))
        rel_path = "./display_config.json"
        abs_file_path = os.path.join(two_up, rel_path)
        abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
        config = open(file=abs_file_path)
        super().__init__(config)
        self.hint = ""

    def get_status(self):
        return {"hint": self.hint}

    def perform_instruction(self, action):
        instruction = action.get("instruction")
        if instruction == "hint":
            self.show_hint(action)
        else:
            return False, action
        return True, None

    def test(self):
        self.hint = "test"
        print(self.hint)
        self.status_changed()

    def show_hint(self, data):
        self.hint = data.get("value")
        print(self.hint)
        self.status_changed()

    def reset(self):
        self.hint = ""
        self.status_changed()

    def main(self):
        self.start()


if __name__ == "__main__":
    device = Display()
    device.main()
