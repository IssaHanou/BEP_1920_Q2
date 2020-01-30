import json
import os
import time

import requests
from sciler.device import Device

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
        rel_path = "hue_lights_config.json"
        abs_file_path = os.path.join(two_up, rel_path)
        abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
        config = open(file=abs_file_path)
        super().__init__(config)
        self.scene = "none"
        self.bri = 0
        self.x = 0
        self.y = 0
        self.hue_bridge = "http://192.168.0.106/"
        self.hue_user = "d3Vji9wgd150ttFBQM3wHl-DyXVBYWnZdO6ALHci"
        self.group = "Spotlights"
        self.header = {"Content-type": "application/json"}

    def get_status(self):
        return {"all": self.scene}

    def perform_instruction(self, action):
        instruction = action.get("instruction")

        if instruction == "scene":
            self.set_scene(action)
        elif instruction == "manual":
            self.set_manual(action.get("component_id"), action.get("value"))
        elif instruction == "bri" or instruction == "x" or instruction == "y":
            self.set_single(action)
        else:
            return False
        return True

    def test(self):
        params = json.dumps({"on": True, "bri": 200, "xy": [0.3, 0.3]})
        url = (self.hue_bridge
               + "api/"
               + self.hue_user
               + "/groups/"
               + self.group
               + "/action")
        requests.put(url,
                     data=params,
                     headers=self.header,
                     )
        time.sleep(2)
        self.pub_to_hue(url)

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
        self.bri = data[1]
        self.x = data[2][0]
        self.y = data[2][1]
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
        self.pub_to_hue(url)

    def set_single(self, action):
        if action.get("instruction") == "bri":
            # self.bri = action.get("value") * 2.5
            self.bri = 100
        elif action.get("instruction") == "x":
            self.x = float(1 / 100 * action.get("value"))
        elif action.get("instruction") == "y":
            self.y = float(1 / 100 * action.get("value"))
        if action.get("component_id") == "all":
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
                    + action.get("component_id")[-1:]
                    + "/state"
            )
        self.pub_to_hue(url)

    def pub_to_hue(self, url):
        params = json.dumps({"on": True, "bri": self.bri, "xy": [self.x, self.y]})
        print(url, params, self.header)
        resp = requests.put(url, data=params, headers=self.header)
        if resp.status_code == 200:
            self.log("Template has been published.")
        else:
            self.log("Unable to publish template.")
        self.status_changed()

    def reset(self):
        self.scene = "none"
        self.bri = 100
        self.x = 0.3
        self.y = 0.3
        url = (
                self.hue_bridge
                + "api/"
                + self.hue_user
                + "/groups/"
                + self.group
                + "/action"
        )
        self.pub_to_hue(url)

    def __loop(self):
        url = (
                self.hue_bridge
                + "api/"
                + self.hue_user
                + "/groups/"
                + self.group
                + "/action"
        )
        while True:
            time.sleep(1)
            self.pub_to_hue(url)

    def main(self):
        self.start()


if __name__ == "__main__":
    device = HueLights()
    device.main()