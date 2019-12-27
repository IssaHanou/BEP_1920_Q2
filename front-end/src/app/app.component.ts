import { Component, OnDestroy, OnInit, ViewEncapsulation } from "@angular/core";
import { IMqttMessage, MqttService } from "ngx-mqtt";
import { Message } from "./message";
import { JsonConvert } from "json2typescript";
import { MatSnackBar, MatSnackBarConfig } from "@angular/material";
import { Subscription } from "rxjs";
import { Devices } from "./components/device/devices";
import { Timers } from "./components/timer/timers";
import { Camera } from "./camera/camera";
import {CameraComponent} from "./camera/camera.component";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.css", "../assets/css/main.css"],
  encapsulation: ViewEncapsulation.None
})
export class AppComponent implements OnInit, OnDestroy {
  title = "S.C.I.L.E.R :";
  nameOfRoom = "Super awesome escape";
  jsonConvert: JsonConvert;
  subscription: Subscription;
  topics = ["front-end"];
  deviceList: Devices;
  timerList: Timers;
  cameras: Camera[];
  selectedCamera: string;
  configErrorList: string[];

  constructor(private mqttService: MqttService, private snackBar: MatSnackBar) {
  }

  /**
   * Initialize app, also called upon loading new config file.
   */
  ngOnInit(): void {
    this.jsonConvert = new JsonConvert();
    this.deviceList = new Devices();
    this.timerList = new Timers();
    this.cameras = [];
    this.configErrorList = [];
    const generaltimer = { id: "general", duration: 0, state: "stateIdle" };
    this.timerList.setTimer(generaltimer);

    for (const topic of this.topics) {
      this.subscribeNewTopic(topic);
    }
    this.sendInstruction([{ instruction: "send status" }]);
    this.sendInstruction([{ instruction: "cameras" }]);
    this.sendConnection(true);
  }

  /**
   * The purpose of this is, when the user leave the app we should cleanup our subscriptions
   * and close the connection with the broker
   */
  ngOnDestroy(): void {
    this.sendConnection(false);
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
   * Send an instruction to the broker, over back-end topic.
   * @param instruction instruction to be sent.
   */
  public sendInstruction(instruction: any[]) {
    let msg = new Message(
      "front-end",
      "instruction",
      new Date(),
      instruction
    );
    let jsonMessage: string = JSON.stringify(this.jsonConvert.serialize(msg));
    this.mqttService.unsafePublish("back-end", jsonMessage);
    for (const inst of instruction) {
      if ("config" in inst) {
        msg.contents = {config: "contents to long to print"};
        jsonMessage = JSON.stringify(this.jsonConvert.serialize(msg));
      }
    }
    console.log(
      "log: sent instruction message: " + jsonMessage
    );
  }

  /**
   * Send an status to the broker, over back-end topic.
   * @param start start status to be sent.
   * @param stop stop status to be sent.
   */
  public sendStatus(start, stop) {
    const msg = new Message("front-end", "status", new Date(), {
      start,
      stop
    });
    const jsonMessage: string = this.jsonConvert.serialize(msg);
    this.mqttService.unsafePublish("back-end", JSON.stringify(jsonMessage));
    console.log("log: sent status message: " + JSON.stringify(jsonMessage));
  }

  /**
   * Send an connection update to the broker, over back-end topic.
   * @param connected connection status to be sent.
   */
  public sendConnection(connected: boolean) {
    const msg = new Message("front-end", "connection", new Date(), {
      connection: connected
    });
    const jsonMessage: string = this.jsonConvert.serialize(msg);
    this.mqttService.unsafePublish("back-end", JSON.stringify(jsonMessage));
    console.log("log: sent connection message: " + JSON.stringify(jsonMessage));
  }

  /**
   * Process incoming message.
   * @param jsonMessage json message.
   */
  public processMessage(jsonMessage: string) {
    const msg: Message = Message.deserialize(jsonMessage);
    switch (msg.type) {
      case "confirmation": {
        this.processConfirmation(jsonMessage);
        break;
      }
      case "instruction": {
        this.processInstruction(msg);
        break;
      }
      case "status": {
        this.deviceList.setDevice(msg.contents);
        break;
      }
      case "time": {
        this.timerList.setTimer(msg.contents);
        break;
      }
      case "cameras": {
        for (const obj of msg.contents) {
          this.cameras.push(new Camera(obj));
        }
        break;
      }
      case "config": {
        this.configErrorList = msg.contents.errors;
        break;
      }
      case "new config": {
        this.ngOnInit();
        break;
      }
      default:
        console.log("log: received invalid message type " + msg.type);
        break;
    }
  }

  /**
   * When the front-end receives confirmation message from client computer
   * that instruction was completed, show the message to the user.
   */
  public processConfirmation(jsonData) {
    const keys = ["instructed", "contents", "instruction"];
    for (const instruction of jsonData.contents[keys[0]][keys[1]]) {
      const display =
        "received confirmation from " +
        jsonData.deviceId +
        " for instruction: " +
        instruction[keys[2]];
      this.openSnackbar(display, "");
    }
  }

  /**
   * Process instruction messages. Two types exist: reset and status update.
   */
  public processInstruction(jsonData) {
    for (const action of jsonData) {
      switch (action.instruction) {
        case "reset":
        {
          this.deviceList.setDevice({
            id: "front-end",
            connection: true,
            status: {
              start: 0,
              stop: 0
            }
          });
        }
          break;
        case "status update": {
          this.sendConnection(true);
        }
      }
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
