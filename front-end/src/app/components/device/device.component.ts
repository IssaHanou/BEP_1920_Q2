import { Component, OnDestroy, OnInit } from "@angular/core";
import { Subscription } from "rxjs";
import { IMqttMessage, MqttService } from "ngx-mqtt";
import { Devices } from "./devices";

@Component({
  selector: "app-device",
  templateUrl: "./device.component.html",
  styleUrls: ["./device.component.css", "../../../assets/css/main.css"]
})
export class DeviceComponent implements OnInit {
  msg: IMqttMessage;
  private subscription: Subscription;
  topicname = "front-end";
  deviceList: Devices;

  onMessageArrived(message) {
    let jsonData = JSON.parse(message);
    jsonData = jsonData.contents;
    this.deviceList.setDevice(jsonData);
  }

  subscribeNewTopic(): void {
    console.log("inside subscribe new topic");
    this.subscription = this.mqttService
      .observe(this.topicname)
      .subscribe((message: IMqttMessage) => {
        this.msg = message;
        console.log(
          "Message: " +
            message.payload.toString() +
            "<br> for topic: " +
            message.topic
        );
        this.onMessageArrived(message.payload);
      });
    console.log("subscribed to topic: " + this.topicname);
  }

  constructor(private mqttService: MqttService) {
    this.subscribeNewTopic();
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
