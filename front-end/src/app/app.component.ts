import { Component, OnDestroy, OnInit, ViewEncapsulation } from "@angular/core";
import { IMqttMessage, MqttService } from "ngx-mqtt";
import { Message } from "./message";
import { JsonConvert } from "json2typescript";
import { MatSnackBar, MatSnackBarConfig } from "@angular/material";
import { Subscription } from "rxjs";
import { Devices } from "./components/device/devices";
import { Puzzles } from "./components/puzzle/puzzles";
import { Timers } from "./components/timer/timers";

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
  puzzleList: Puzzles;
  timerList: Timers;

  constructor(private mqttService: MqttService, private snackBar: MatSnackBar) {
    this.jsonConvert = new JsonConvert();
    this.deviceList = new Devices();
    this.puzzleList = new Puzzles();
    this.timerList = new Timers();
    const generalTimer = { id: "general", duration: 0, state: "stateIdle" };
    this.timerList.setTimer(generalTimer);
  }

  ngOnInit(): void {
    for (const topic of this.topics) {
      this.subscribeNewTopic(topic);
    }
    this.sendInstruction([{ instruction: "name" }]);
    this.sendInstruction([{ instruction: "send status" }]);
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
    const msg = new Message(
      "front-end",
      "instruction",
      new Date(),
      instruction
    );
    const jsonMessage: string = this.jsonConvert.serialize(msg);
    this.mqttService.unsafePublish("back-end", JSON.stringify(jsonMessage));
    console.log(
      "log: sent instruction message: " + JSON.stringify(jsonMessage)
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
        /**
         * When the front-end receives confirmation message from client computer
         * that instruction was completed, show the message to the user.
         */

        for (const instruction of msg.contents.instructed.contents) {
          const display =
            "received confirmation from " +
            msg.deviceId +
            " for instruction: " +
            instruction.instruction;
          this.openSnackbar(display, "");
        }
        break;
      }
      case "event status": {
        this.puzzleList.updatePuzzles(msg.contents);
        break;
      }
      case "instruction": {
        for (const action of msg.contents) {
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
        break;
      }
      case "status": {
        this.deviceList.setDevice(msg.contents);
        break;
      }
      case "time": {
        this.processTimeStatus(msg.contents);
        break;
      }
      case "name": {
        this.nameOfRoom = msg.contents.name;
        break;
      }
      default:
        console.log("log: received invalid message type " + msg.type);
        break;
    }
  }

  /**
   * Timers send their status to the front-end but we only care about the general time.
   * @param jsonData with id, status and state
   */
  public processTimeStatus(jsonData) {
    this.timerList.setTimer(jsonData);
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
