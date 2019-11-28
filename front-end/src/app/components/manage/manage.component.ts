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
    let message;
    let now;
    now = new Date();
    message = {
      type: "instruction",
      instruction: "test all",
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
        now.getSeconds()
    };
    this.mqttService.unsafePublish("back-end", JSON.stringify(message));
  }
}
