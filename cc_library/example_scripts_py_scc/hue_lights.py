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


class Spot():

    def __init__(self, x=0.3, y=0.3, bri=100):
        self.x = x
        self.y = y
        self.bri = bri


class HueLights(Device):
    def __init__(self):
        two_up = os.path.abspath(os.path.join(__file__, ".."))
        rel_path = "hue_lights_config.json"
        abs_file_path = os.path.join(two_up, rel_path)
        abs_file_path = os.path.abspath(os.path.realpath(abs_file_path))
        config = open(file=abs_file_path)
        super().__init__(config)
        self.scene = "none"
        self.spots = [Spot(), Spot(), Spot(), Spot()]
        self.hue_bridge = "http://192.168.0.107/"
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
        bri = data[1]
        x = data[2][0]
        y = data[2][1]
        if comp == "all":
            for i in range(len(self.spots)):
                self.set_spot(i + 1, x, y, bri)
        else:
            self.set_spot(comp[-1:], x, y, bri)

    def set_spot(self, spot, x=None, y=None, bri=None):
        if x:
            self.spots[spot].x = x
        if y:
            self.spots[spot].y = y
        if bri:
            self.spots[spot].bri = bri

    def set_single(self, action):
        component = action.get("component")
        if component == "all":
            for i in range(len(self.spots)):
                self.set_single_spot(i + 1, action)
        else:
            self.set_spot(component[-1:], action)

    def set_single_spot(self, spot, action):
        if action.get("instruction") == "bri":
            bri = int(action.get("value") * 2.54)
            self.set_spot(spot, bri=bri)
        elif action.get("instruction") == "x":
            x = float(1 / 100 * action.get("value"))
            self.set_spot(spot, x=x)
        elif action.get("instruction") == "y":
            y = float(1 / 100 * action.get("value"))
            self.set_spot(spot, y=y)

    def check_differences_and_publish(self, previous):
        url = None
        params = None
        head = self.spots[0]
        allEqual = True
        for spot in self.spots:
            if spot != head:
                allEqual = False

        if allEqual:
            url = (
                    self.hue_bridge
                    + "api/"
                    + self.hue_user
                    + "/groups/"
                    + self.group
                    + "/action"
            )
            params = json.dumps({"on": True, "bri": head.bri, "xy": [head.x, head.y], "transitiontime": 5})
            self.pub_to_hue(url, params)
        else:
            for i in range(len(self.spots)):
                current = self.spots[i]
                if current != previous[i]:
                    url = (
                            self.hue_bridge
                            + "api/"
                            + self.hue_user
                            + "/lights/"
                            + str(i+1)
                            + "/state"
                    )
                    params = json.dumps({"on": True, "bri": current.bri, "xy": [current.x, current.y], "transitiontime": 5})
                    self.pub_to_hue(url, params)


    def pub_to_hue(self, url, params):
        print(url, params, self.header)
        resp = requests.put(url, data=params, headers=self.header)
        if resp.status_code == 200:
            self.log("Template has been published.")
        else:
            self.log("Unable to publish template.")
        self.status_changed()


    def reset(self):
        self.scene = "none"
        self.spots = [Spot(), Spot(), Spot(), Spot()]

    def __loop(self):
        previous = self.spots
        while True:
            if self.spots != previous:
                self.check_differences_and_publish(previous)
                previous = self.spots
                time.sleep(1)

    def main(self):
        self.start(self.__loop)


if __name__ == "__main__":
    device = HueLights()
    device.main()
