import { Component, OnInit } from "@angular/core";
import { AppComponent } from "../../app.component";
import { formatTime } from "./timer";

/**
 * The timer component is shown on the front-page in the "Tijd" box.
 */
@Component({
  selector: "app-timer",
  templateUrl: "./timer.component.html",
  styleUrls: ["./timer.component.css", "../../../assets/css/main.css"]
})
export class TimerComponent implements OnInit {
  constructor(private app: AppComponent) {}

  ngOnInit() {}

  getDisplayTime() {
    return this.app.displayTime;
  }

  /**
   * The done time is displayed when the game has started.
   * It shows when the game will finish, depending on the current time and the duration of the game.
   */
  getDoneTime() {
    const device = this.app.deviceList.getDevice("front-end");
    if (device !== null) {
      const status = device.status;
      if (status.get("start") > 0 && status.get("stop") === 0) {
        const date = new Date();
        const timeDiff =
          date.getTime() + this.app.timerList.getTimer("general").duration;
        const doneTime = formatTime(timeDiff, date.getTimezoneOffset());
        return "Klaar om " + doneTime;
      }
    }
  }
}
