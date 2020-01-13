import json
import os
import time

import requests
from scclib.device import Device

"""
How to use Hue Lights with S.C.I.L.E.R.:
- download Hue app
- create groups and scenes to use in the escape room
- retrieve IP address of hue bridge ( apt-get install avahi-utils ; avahi-browse -rt _hue._tcp )
- retrieve valid hue username (curl -d '{"devicetype":"["whatever"]"}' -H
"Content-Type: application/json" -X POST 'http://<BRIDGE_IP>/api' ; returns long hue username)
- retrieve scene id's ( curl 'http://<hue bridge IP>/api/<hue username>/scenes )
- implement config using scene ids
"""


class HueLights(Device):
    def __init__(self):
        two_up = os.path.abspath(os.path.join(__file__, ".."))
        rel_path = "./hue_lights_config.json"
        abs_file_path = os.path.join(two_up, rel_path)
        abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
        config = open(file=abs_file_path)
        super().__init__(config)
        self.scene = "none"
        self.hue_bridge = "http://192.168.178.128/"
        self.hue_user = "JQrPwJNthHtfPEG9vhW3mqwIVuFo3ESLD3gvkZOB"
        self.group = "Spotlights"
        self.header = {"Content-type": "application/json"}

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
        params = json.dumps({"on": True, "bri": 200, "xy": [0.3, 0.3]})
        requests.put(
            self.hue_bridge
            + "api/"
            + self.hue_user
            + "/groups/"
            + self.group
            + "/action",
            data=params,
            headers=self.header,
        )
        time.sleep(2)
        params = json.dumps({"scene": self.scene})
        resp = requests.put(
            self.hue_bridge
            + "api/"
            + self.hue_user
            + "/groups/"
            + self.group
            + "/action",
            data=params,
            headers=self.header,
        )
        if resp.status_code == 200:
            self.log("Template has been published.")
        else:
            self.log("Unable to publish template.")
        self.status_changed()

    def set_scene(self, data):
        self.scene = data.get("value")
        params = json.dumps({"scene": self.scene})
        resp = requests.put(
            self.hue_bridge
            + "api/"
            + self.hue_user
            + "/groups/"
            + self.group
            + "/action",
            data=params,
            headers=self.header,
        )
        if resp.status_code == 200:
            self.log("Template has been published.")
        else:
            self.log("Unable to publish template.")
        self.status_changed()

    def set_manual(self, comp, data):
        params = json.dumps({"on": data[0], "bri": data[1], "xy": data[2]})
        if comp == "all":
            url = (
                self.hue_bridge
                + "api/"
                + self.hue_user
                + "/groups/"
                + self.group
                + "/action"
            )
        else:
            url = (
                self.hue_bridge
                + "api/"
                + self.hue_user
                + "/lights/"
                + comp[-1:]
                + "/state"
            )
        resp = requests.put(url, data=params, headers=self.header)
        if resp.status_code == 200:
            self.log("Template has been published.")
        else:
            self.log("Unable to publish template.")
        self.status_changed()

    def reset(self):
        self.scene = "none"
        params = json.dumps({"on": True, "bri": 50, "xy": [0.3, 0.3]})
        resp = requests.put(
            self.hue_bridge
            + "api/"
            + self.hue_user
            + "/groups/"
            + self.group
            + "/action",
            data=params,
            headers=self.header,
        )
        if resp.status_code == 200:
            self.log("action has been published.")
        else:
            self.log("Unable to publish template.")
        self.status_changed()

    def main(self):
        self.start()


if __name__ == "__main__":
    device = HueLights()
    device.main()
