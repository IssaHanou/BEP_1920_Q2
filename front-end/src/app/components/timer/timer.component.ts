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
