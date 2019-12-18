import { Component, OnInit } from "@angular/core";
import * as moment from "moment";
import { Subscription, Observable, timer } from "rxjs";
import { AppComponent } from "../../app.component";

@Component({
  selector: "app-timer",
  templateUrl: "./timer.component.html",
  styleUrls: ["./timer.component.css", "../../../assets/css/main.css"]
})
export class TimerComponent implements OnInit {
  private subscription: Subscription;
  displayTime: string;
  everySecond: Observable<number> = timer(0, 1000);

  constructor(private app: AppComponent) {}

  ngOnInit() {
    this.subscription = this.everySecond.subscribe(seconds => {
      if (this.app.timeState === "stateActive") {
        this.app.remainingTime = this.app.remainingTime - 1000;
        this.displayTime = formatMS(this.app.remainingTime);
      } else if (this.app.timeState === "stateIdle") {
        this.displayTime = formatMS(this.app.remainingTime);
      }
    });
  }
}

export function formatMS(timeInMS) {
  const seconds = parseInt(((timeInMS / 1000) % 60).toString(), 10);
  const minutes = parseInt(((timeInMS / (1000 * 60)) % 60).toString(), 10);
  const hours = parseInt(((timeInMS / (1000 * 60 * 60)) % 24).toString(), 10);
  const h = hours < 10 ? "0" + hours : hours;
  const m = minutes < 10 ? "0" + minutes : minutes;
  const s = seconds < 10 ? "0" + seconds : seconds;

  return h + ":" + m + ":" + s;
  // return moment(timeInMS).format("hh:mm:ss"); // TODO: Adds an hour
}
