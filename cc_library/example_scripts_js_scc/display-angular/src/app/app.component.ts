import { Component, OnInit } from "@angular/core";
import { Display } from "./display";
import * as data from "../assets/display_config.json";
import { Observable, Subscription, timer } from "rxjs";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.css"]
})
export class AppComponent implements OnInit {
  sub: Subscription;

  title = "display";
  config;
  display;
  displayTime;
  everySecond: Observable<number> = timer(0, 1000);

  constructor() {
    this.config = (data as any).default;
    this.display = new Display(this.config);
  }

  ngOnInit(): void {
    this.display.start(() => {
      console.log("CONNECTED");
    });
    this.initializeTimers();
  }

  /**
   * Initialize the timer to listen to everySecond and set their state accordingly.
   */
  private initializeTimers() {
    this.sub = this.everySecond.subscribe(seconds => {
      if (this.display.timeState === "stateActive") {
        this.display.timeDur = this.display.timeDur - 1000;
      }
      if (this.display.timeDur <= 0) {
        this.display.timeState = "stateIdle";
      }
      this.displayTime = this.formatMS(this.display.timeDur);
    });
  }

  /**
   * Format the time in milliseconds to a string in the format hh:mm:ss.
   */
  formatMS(timeInMS) {
    const seconds = parseInt(((timeInMS / 1000) % 60).toString(), 10);
    const minutes = parseInt(((timeInMS / (1000 * 60)) % 60).toString(), 10);
    const hours = parseInt(((timeInMS / (1000 * 60 * 60)) % 24).toString(), 10);
    const h = hours < 10 ? "0" + hours : hours;
    const m = minutes < 10 ? "0" + minutes : minutes;
    const s = seconds < 10 ? "0" + seconds : seconds;

    return h + ":" + m + ":" + s;
  }
}
