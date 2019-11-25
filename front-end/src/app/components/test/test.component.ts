import { Component, OnInit } from "@angular/core";
import { MqttService } from "ngx-mqtt";

@Component({
  selector: "app-test",
  templateUrl: "./test.component.html",
  styleUrls: ["./test.component.css"]
})
export class TestComponent {
  constructor(private mqttService: MqttService) {}

  onClick() {
    let message;
    let now;
    now = new Date();
    message = {
      type: "instruction",
      instruction: "test all",
      time_sent:
        now.getDay() +
        "-" +
        now.getMonth() +
        "-" +
        now.getFullYear() +
        " " +
        now.getHours() +
        ":" +
        now.getMinutes() +
        ":" +
        now.getSeconds()
    };
    this.mqttService.unsafePublish("back-end", JSON.stringify(message));
  }
}
