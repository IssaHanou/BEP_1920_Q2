import { Component, OnInit } from "@angular/core";
import { AppComponent } from "../../app.component";

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
}

export function formatMS(timeInMS) {
  const seconds = parseInt(((timeInMS / 1000) % 60).toString(), 10);
  const minutes = parseInt(((timeInMS / (1000 * 60)) % 60).toString(), 10);
  const hours = parseInt(((timeInMS / (1000 * 60 * 60)) % 24).toString(), 10);
  const h = hours < 10 ? "0" + hours : hours;
  const m = minutes < 10 ? "0" + minutes : minutes;
  const s = seconds < 10 ? "0" + seconds : seconds;

  return h + ":" + m + ":" + s;
}
