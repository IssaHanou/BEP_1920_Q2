import { Component, OnInit } from '@angular/core';
import * as moment from "moment";

@Component({
  selector: 'timer',
  templateUrl: './timer.component.html',
  styleUrls: ['./timer.component.css', '../../../assets/css/main.css']
})
export class TimerComponent implements OnInit {

  remaining_time = formatMS(179000);
  constructor() { }

  ngOnInit() {
  }
}

export function formatMS(timeInMS) {
  return moment(timeInMS).format("hh:mm:ss"); //Adds an hour
}
