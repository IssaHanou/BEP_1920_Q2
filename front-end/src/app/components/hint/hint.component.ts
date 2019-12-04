import { Component } from "@angular/core";
import { MqttService } from "ngx-mqtt";

@Component({
  selector: "app-hint",
  templateUrl: "./hint.component.html",
  styleUrls: ["./hint.component.css", "../../../assets/css/main.css"]
})
export class HintComponent {
  hint: string;

  constructor(private mqttService: MqttService) {
  }

  public unsafePublish(topic: string, message: string): void {
    this.mqttService.unsafePublish(topic, message, { qos: 1, retain: true });
  }

  onSubmit() {
    this.unsafePublish("hint", this.hint);
    this.hint = "";
  }
}
