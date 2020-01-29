import { Component, OnInit } from "@angular/core";
import { AppComponent } from "../../app.component";
import { EventsMethods } from "../EventsMethods";

@Component({
  selector: "app-event",
  templateUrl: "./event.component.html",
  styleUrls: ["./event.component.css", "../../../assets/css/main.css", "./../events.css"],
  providers: [EventsMethods]
})
export class EventComponent implements OnInit {

  constructor(private app: AppComponent, private methods: EventsMethods) {}

  ngOnInit() {
  }
}
