import { Component, OnInit } from "@angular/core";
import * as data from "../assets/display_config.json";
const SccLib = require("../../../../js-scc"); // development
// const SccLib = require("js-scc"); // production

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.css"]
})
export class AppComponent implements OnInit {
  title = "display";
  hint = "";
  config;
  scc;

  constructor() {
    this.config = (data as any).default;
    this.scc = new SccLib(this.config, 4, function(date, level, message) {
      const formatDate = function(date) {
        return (
          date.getDate() +
          "-" +
          date.getMonth() +
          1 +
          "-" +
          date.getFullYear() +
          " " +
          date.getHours() +
          ":" +
          date.getMinutes() +
          ":" +
          date.getSeconds()
        );
      };
      console.log(
        "time=" + formatDate(date) + " level=" + level + " msg=" + message
      ); // call own logger
    });
  }

  ngOnInit(): void {
    this.scc.connect();
  }
}
