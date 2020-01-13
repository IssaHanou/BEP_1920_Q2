import { Component, OnInit } from "@angular/core";
import { Display } from "./display";
import * as data from "../assets/display_config.json";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.css"]
})
export class AppComponent implements OnInit {
  title = "display";
  config;
  display;

  constructor() {
    this.config = (data as any).default;
    this.display = new Display(this.config);
  }

  ngOnInit(): void {
    this.display.start();
  }
}
