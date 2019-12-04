import { Component, OnDestroy, OnInit } from "@angular/core";
import {Observable} from 'rxjs/internal/Observable';

@Component({
  selector: "app-device",
  templateUrl: "./device.component.html",
  styleUrls: ["./device.component.css", "../../../assets/css/main.css"]
})


export class DeviceComponent implements OnInit {
  constructor() {
  }

  data
  columns

  ngOnInit() {
    this.columns = ["Apparaat", "Connectie", "Status"]
    this.data = [
      {
        Apparaat: "Weegschaal",
        Connectie: "true",
        Status: "40",
      },
      {
        Apparaat: "telefoon",
        Connectie: "true",
        Status: "0612254049",
      },
      {
        Apparaat: "Controlbord",
        Connectie: "false",
        Status: "false" ,
      },
    ]
  }
}
