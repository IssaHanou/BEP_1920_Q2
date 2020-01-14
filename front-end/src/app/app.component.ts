import { Component, OnDestroy, OnInit, ViewEncapsulation } from "@angular/core";
import { IMqttMessage, MqttService } from "ngx-mqtt";
import { Message } from "./message";
import { JsonConvert } from "json2typescript";
import { MatSnackBar, MatSnackBarConfig } from "@angular/material";
import { Observable, Subscription, timer } from "rxjs";
import { Devices } from "./components/device/devices";
import { Puzzles } from "./components/puzzle/puzzles";
import { Timers } from "./components/timer/timers";
import { Logger } from "./logger";
import { Camera } from "./camera/camera";
import { Hint } from "./components/hint/hint";
import { formatMS, formatTime } from "./components/timer/timer";
import { FullScreen } from "./fullscreen";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.css", "../assets/css/main.css"],
  encapsulation: ViewEncapsulation.None
})
export class AppComponent extends FullScreen implements OnInit, OnDestroy {
  // Variables for the home screen
  title = "SCILER";
  nameOfRoom = "Super awesome escape";

  // Necessary tools
  jsonConvert: JsonConvert;
  logger: Logger;
  subscription: Subscription;

  // Keeping track of data
  deviceList: Devices;
  puzzleList: Puzzles;
  hintList: Hint[];
  configErrorList: string[];
  cameras: Camera[];
  selectedCamera: string;
  selectedCamera2: string;
  openSecondCamera = false;
  timerList: Timers;
  displayTime: string;
  everySecond: Observable<number> = timer(0, 1000);

  constructor(private mqttService: MqttService, private snackBar: MatSnackBar) {
    super();
    this.logger = new Logger();
    this.jsonConvert = new JsonConvert();
    this.deviceList = new Devices();
    this.puzzleList = new Puzzles();
    this.hintList = [];
    this.configErrorList = [];
    this.cameras = [];

    const topics = ["front-end"];
    for (const topic of topics) {
      this.subscribeNewTopic(topic);
    }
  }

  /**
   * Initialize app, also called upon loading new config file.
   */
  ngOnInit(): void {
    this.timerList = new Timers();
    const generalTimer = { id: "general", duration: 0, state: "stateIdle" };
    this.timerList.setTimer(generalTimer);

    this.sendInstruction([{ instruction: "send setup" }]);
    this.sendConnection(true);
    this.initializeTimers();
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
  private subscribeNewTopic(topic: string): void {
    this.subscription = this.mqttService
      .observe(topic)
      .subscribe((message: IMqttMessage) => {
        this.logger.log("info",
          "received on topic " +
            message.topic +
            ", message: " +
            message.payload.toString()
        );
        this.processMessage(message.payload.toString());
      });
    this.logger.log("info", "subscribed to topic: " + topic);
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
    let jsonMessage: string = JSON.stringify(this.jsonConvert.serialize(msg));
    this.mqttService.unsafePublish("back-end", jsonMessage);
    for (const inst of instruction) {
      if ("config" in inst) {
        msg.contents = {config: "contents to long to print"};
        jsonMessage = JSON.stringify(this.jsonConvert.serialize(msg));
      }
    }
    this.logger.log("info",
      "sent instruction message: " + jsonMessage
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
    this.logger.log("info", "sent status message: " + JSON.stringify(jsonMessage));
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
    this.logger.log("info", "sent connection message: " + JSON.stringify(jsonMessage));
  }

  /**
   * Process incoming message.
   * @param jsonMessage json message.
   */
  private processMessage(jsonMessage: string) {
    const msg: Message = Message.deserialize(jsonMessage);
    switch (msg.type) {
      case "confirmation": {
        this.processConfirmation(msg);
        break;
      }
      case "instruction": {
        this.processInstruction(msg.contents);
        break;
      }
      case "status": {
        this.deviceList.setDevice(msg.contents);
        break;
      }
      case "event status": {
        this.puzzleList.updatePuzzles(msg.contents);
        break;
      }
      case "time": {
        this.timerList.setTimer(msg.contents);
        break;
      }
      case "setup": {
        this.processSetUp(msg.contents);
        break;
      }
      case "config": {
        this.configErrorList = msg.contents.errors;
        break;
      }
      case "new config": {
        this.stopTimers();
        this.ngOnInit();
        this.openSnackbar("using new config: " + msg.contents.name, "");
        break;
      }
      default:
        this.logger.log("error", "received invalid message type " + msg.type);
        break;
    }
  }

  /**
   * When the front-end receives confirmation message from client computer
   * that instruction was completed, show the message to the user.
   */
  private processConfirmation(jsonData) {
    for (const instruction of jsonData.contents.instructed.contents) {
      const display =
        "received confirmation from " +
        jsonData.deviceId +
        " for instruction: " +
        instruction.instruction;
      this.openSnackbar(display, "");
    }
  }

  /**
   * Process instruction messages. Two types exist: reset and status update.
   */
  private processInstruction(jsonData) {
    for (const action of jsonData) {
      switch (action.instruction) {
        case "reset": {
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
          break;
        }
        case "test": {
          this.openSnackbar("performing instruction test", "");
          break;
        }
      }
    }
  }

  /**
   * The setup contains the name of the room, the map with hints per puzzle and the rule descriptions.
   * @param jsonData with name, hints, events
   */
  private processSetUp(jsonData) {
    this.nameOfRoom = jsonData.name;

    const cameraData = jsonData.cameras;
    this.cameras = [];
    if (cameraData !== null) {
      for (const cam of cameraData) {
        this.cameras.push(new Camera(cam));
      }
    }

    const rules = jsonData.events;
    this.puzzleList = new Puzzles();
    for (const rule in rules) {
      if (rules.hasOwnProperty(rule)) {
        this.puzzleList.addPuzzle(rule, rules[rule]);
      }
    }

    const allHints = jsonData.hints;
    this.hintList = [];
    for (const puzzle in allHints) {
      if (allHints.hasOwnProperty(puzzle)) {
        const hints = [];
        for (const index in allHints[puzzle]) {
          if (allHints[puzzle].hasOwnProperty(index)) {
            hints.push(allHints[puzzle][index]);
          }
        }
        this.hintList.push(new Hint(puzzle, hints));
      }
    }
  }

  /**
   * Initialize the timers to listen to every second and set their state accordingly.
   */
  private initializeTimers() {
    this.subscription = this.everySecond.subscribe(seconds => {
      for (const aTimer of this.timerList.getAll().values()) {
        if (aTimer.state === "stateActive") {
          aTimer.tick();
        }
        if (aTimer.duration <= 0) {
          aTimer.state = "stateIdle";
        }
      }
      this.displayTime = formatMS(
        this.timerList.getTimer("general").getTimeLeft()
      );
    });
  }

  /**
   * Before using new configuration, first stop the current timer subscription.
   * Otherwise time runs double.
   */
  private stopTimers() {
    this.subscription.unsubscribe();
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

  /**
   * Return the current time to display.
   */
  getCurrentTime() {
    return formatTime(new Date());
  }
}
