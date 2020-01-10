import json
import os
import time
from scclib.device import Device

import requests

class Hue_lights(Device):
    def __init__(self):
        two_up = os.path.abspath(os.path.join(__file__, ".."))
        rel_path = "./hue_lights_config.json"
        abs_file_path = os.path.join(two_up, rel_path)
        abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
        config = open(file=abs_file_path)
        super().__init__(config)
        self.scene = "none"

    def get_status(self):
        return {"all": self.scene}

    def perform_instruction(self, action):
        instruction = action.get("instruction")
        if instruction == "scene":
            self.set_scene(action)
        if instruction == "manual":
            self.set_manual(action.get("component_id"), action.get("value"))
        else:
            return False, action
        return True, None

    def test(self):
        print("test")
        params = json.dumps({'on': True, 'bri': 200, "xy": [0.3,0.3]})
        headers = {"Content-type": "application/json"}
        requests.put("http://192.168.178.128/api/JQrPwJNthHtfPEG9vhW3mqwIVuFo3ESLD3gvkZOB/groups/Spotlights/action", data=params, headers=headers)
        time.sleep(2)
        params = json.dumps({"scene": self.scene})
        headers = {"Content-type": "application/json"}
        resp = requests.put("http://192.168.178.128/api/JQrPwJNthHtfPEG9vhW3mqwIVuFo3ESLD3gvkZOB/groups/Spotlights/action", data=params, headers=headers)
        if resp.status_code == 200:
            print('Template has been published.')
        else:
            print('Unable to publish template.')
        self.status_changed()

    def set_scene(self, data):
        self.scene = data.get("value")
        params = json.dumps({"scene": self.scene})
        headers = {"Content-type": "application/json"}
        resp = requests.put("http://192.168.178.128/api/JQrPwJNthHtfPEG9vhW3mqwIVuFo3ESLD3gvkZOB/groups/Spotlights/action", data=params, headers=headers)
        if resp.status_code == 200:
            print('Template has been published.')
        else:
            print('Unable to publish template.')
        self.status_changed()

    def set_manual(self, comp, data):
        params = json.dumps({'on': data[0], 'bri': data[1], 'xy':data[2]})
        headers = {"Content-type": "application/json"}
        url = "http://192.168.178.128/api/JQrPwJNthHtfPEG9vhW3mqwIVuFo3ESLD3gvkZOB/lights/" + comp[-1:] + "/state"
        print(url, params, comp[-1:])
        resp = requests.put(url, data=params, headers=headers)
        if resp.status_code == 200:
            print('Template has been published.')
        else:
            print('Unable to publish template.')
        self.status_changed()

    def reset(self):
        self.scene = "none"
        params = json.dumps({'on': True, 'bri': 50, "xy": [0.3,0.3]})
        headers = {"Content-type": "application/json"}
        resp = requests.put(
            "http://192.168.178.128/api/JQrPwJNthHtfPEG9vhW3mqwIVuFo3ESLD3gvkZOB/groups/Spotlights/action", data=params,
            headers=headers)
        if resp.status_code == 200:
            print('Template has been published.')
        else:
            print('Unable to publish template.')
        self.status_changed()

    def main(self):
        self.start()


if __name__ == "__main__":
    device = Hue_lights()
    device.main()
