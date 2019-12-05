import {Component, OnDestroy, OnInit} from "@angular/core";
import {Subscription} from "rxjs";
import {IMqttMessage, MqttService} from "ngx-mqtt";

@Component({
  selector: "app-device",
  templateUrl: "./device.component.html",
  styleUrls: ["./device.component.css", "../../../assets/css/main.css"]
})


export class DeviceComponent implements OnInit {

  msg: IMqttMessage;
  private subscription: Subscription;
  topicname= "front-end";
  deviceList: Devices;

  onMessageArrived(message) {
    let jsonData  = JSON.parse(message);
    jsonData = jsonData.contents;
    this.deviceList.setDevice(jsonData);
  }

  subscribeNewTopic(): void {
    console.log("inside subscribe new topic");
    this.subscription = this.mqttService.observe(this.topicname).subscribe((message: IMqttMessage) => {
      this.msg = message;
      console.log("Message: " + message.payload.toString() + "<br> for topic: " + message.topic);
      this.onMessageArrived(message.payload)
    });
    console.log("subscribed to topic: " + this.topicname)
  }


  constructor(private mqttService: MqttService) {
    this.subscribeNewTopic()
  }

  ngOnInit() {
    this.deviceList = new Devices();
    const now = new Date();
    const message = {
      device_id: "front-end",
      time_sent:
        now.getDate() +
        "-" +
        (now.getMonth() + 1) +
        "-" +
        now.getFullYear() +
        " " +
        now.getHours() +
        ":" +
        now.getMinutes() +
        ":" +
        now.getSeconds(),
      type: "instruction",
      contents: {
        instruction: "send status"
      }
    };
    this.mqttService.unsafePublish("back-end", JSON.stringify(message));
  }
}

/**
 * Devices has a Map all containing all devices with a key that is the same as the id.
 */
export class Devices {
  all: Map<string, Device>;

  constructor() {
    this.all = new Map<string,Device>()
  }

  /**
   * setDevice either updates an existing Device with the update methods or creates a new one.
   * @param jsonData json object with keys id, status and connection. 
   */
  setDevice(jsonData){
    if (this.all.has(jsonData.id)) {
      this.all.get(jsonData.id).updateStatus(jsonData.status);
      this.all.get(jsonData.id).updateConnection(jsonData.connection)
    } else {
      this.all.set(jsonData.id, new Device(jsonData))
    }

  }
}

/**
 * Device has all the information the front-end needs to show of each device
 */
export class Device {
  id: string;
  status: Map<string, any>;
  connection: boolean;

  constructor(jsonData) {
    this.id = jsonData.id;
    this.updateConnection(jsonData.connection);
    this.status = new Map<string, any>();
    this.updateStatus(jsonData.status)
  }

  /**
   * updateConnection is called on every status update to update the connections status.
   * @param connection boolean
   */
  updateConnection(connection){
    this.connection = connection;
  }

  /**
   * updateStatus is called on every status update to update the component status.
   * @param jsonStatus json object containing components as key with status as value.
   */
  updateStatus(jsonStatus) {
    const keys = Object.keys(jsonStatus);
    for (const key of keys) {
      if (this.status.has(key)) {
        this.status.delete(key);
        this.status.set(key, jsonStatus[key])
      } else {
        this.status.set(key, jsonStatus[key])
      }
    }
  }

}

