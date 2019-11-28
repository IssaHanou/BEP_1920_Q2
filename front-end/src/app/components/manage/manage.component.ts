import { Component, OnInit } from "@angular/core";
import { MqttService } from "ngx-mqtt";

@Component({
  selector: "app-manage",
  templateUrl: "./manage.component.html",
  styleUrls: ["./manage.component.css", "../../../assets/css/main.css"]
})
export class ManageComponent implements OnInit {
  constructor(private mqttService: MqttService) {}

  ngOnInit() {}

  onClickTest() {
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
        instruction: "test all"
      }
    };
    this.mqttService.unsafePublish("back-end", JSON.stringify(message));
  }
}
