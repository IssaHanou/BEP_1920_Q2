import { Component, OnDestroy, OnInit, ViewEncapsulation } from "@angular/core";
import { IMqttMessage, MqttService } from "ngx-mqtt";
import { Message } from "./message";
import { JsonConvert } from "json2typescript";
import { MatSnackBar, MatSnackBarConfig } from "@angular/material";
import { Subscription } from "rxjs";
import { Devices } from "./components/device/devices";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: [
    "./app.component.css",
    "../assets/css/main.css"
  ],
  encapsulation: ViewEncapsulation.None
})
export class AppComponent implements OnInit, OnDestroy {
  title = "S.C.I.L.E.R";
  nameOfRoom = "Super awesome escape";
  jsonConvert: JsonConvert;
  subscription: Subscription;
  topics = ["front-end"];
  deviceList: Devices;

  constructor(private mqttService: MqttService, private snackBar: MatSnackBar) {
    this.jsonConvert = new JsonConvert();
    this.deviceList = new Devices();
  }

  ngOnInit(): void {
    for (const topic of this.topics) {
      this.subscribeNewTopic(topic);
    }
    this.sendInstruction([{instruction: "send status"}]);
  }

  /**
   * The purpose of this is, when the user leave the app we should cleanup our subscriptions
   * and close the connection with the broker
   */
  ngOnDestroy(): void {
    this.mqttService.disconnect();
  }

  /**
   * Subscribe to topics.
   */
  public subscribeNewTopic(topic: string): void {
    this.subscription = this.mqttService
      .observe(topic)
      .subscribe((message: IMqttMessage) => {
        console.log(
          "log: received on topic " +
            message.topic +
            ", message: " +
            message.payload.toString()
        );
        this.processMessage(message.payload.toString());
      });
    console.log("log: subscribed to topic: " + topic);
  }

  /**
   * Send an instruction to the broker, over instruction topic.
   * @param instruction instruction to be sent.
   */
  public sendInstruction(instruction: any[]) {
    const msg = new Message("front-end", "instruction", new Date(), instruction);
    const jsonMessage: string = this.jsonConvert.serialize(msg);
    this.mqttService.unsafePublish("back-end", JSON.stringify(jsonMessage));
    console.log(
      "log: sent instruction message: " + JSON.stringify(jsonMessage)
    );
  }

  /**
   * Process incoming message.
   * @param jsonMessage json message.
   */
  public processMessage(jsonMessage: string) {
    const msg: Message = Message.deserialize(jsonMessage);
    switch (msg.type) {
      case "confirmation": {
        const keys = ["instructed", "contents", "instruction"];
        /**
         * When the front-end receives confirmation message from client computer
         * that instruction was completed, show the message to the user.
         */

        for (const instruction of msg.contents[keys[0]][keys[1]]) {
          const display =
            "received confirmation from " +
            msg.deviceId +
            " for instruction: " +
            instruction[keys[2]];
          this.openSnackbar(display, "");
        }
        break;
      }
      case "instruction": {
        // TODO instructions to front-end? e.g. ask for hint
        break;
      }
      case "status": {
        this.deviceList.setDevice(msg.contents);
        break;
      }
      default:
        console.log("log: received invalid message type " + msg.type);
        break;
    }
  }

  /**
   * Opens snackbar with duration of 2 seconds.
   * @param message displays this message
   * @param action: button to display
   */
  public openSnackbar(message: string, action: string) {
    const config = new MatSnackBarConfig();
    config.duration = 3000;
    config.panelClass = ["custom-snack-bar"];
    this.snackBar.open(message, action, config);
  }
}
