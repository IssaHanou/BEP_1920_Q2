import { Component } from "@angular/core";
import { MqttService } from "ngx-mqtt";
import { Message } from "../../message";
import { AppComponent } from "../../app.component";

@Component({
  selector: "app-hint",
  templateUrl: "./hint.component.html",
  styleUrls: ["./hint.component.css", "../../../assets/css/main.css"]
})
export class HintComponent {
  hint: string;

  constructor(private mqttService: MqttService, private app: AppComponent) {}

  public unsafePublish(topic: string, message: string): void {
    this.mqttService.unsafePublish(topic, message, { qos: 1, retain: true });
  }

  onSubmit() {
    const msg = new Message("front-end", "instruction", new Date(), [
      {
        instruction: "hint",
        value: this.hint
      }
    ]);
    const res = this.app.jsonConvert.serialize(msg);
    this.unsafePublish("back-end", JSON.stringify(res));
    console.log("log: sent instruction message: " + JSON.stringify(res));
    this.hint = "";
  }
}
