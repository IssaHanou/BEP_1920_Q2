import { Component, OnInit } from "@angular/core";
import { Device } from "./device";
import { AppComponent } from "../../app.component";

@Component({
  selector: "app-device",
  templateUrl: "./device.component.html",
  styleUrls: ["./device.component.css", "../../../assets/css/main.css"]
})
export class DeviceComponent implements OnInit {
  constructor(private app: AppComponent) {}

  ngOnInit() {}

  public getDeviceStatus(): Device[] {
    const devices: Device[] = [];
    for (const device of this.app.deviceList.all.values()) {
      devices.push(device);
    }
    devices.sort((a: Device, b: Device) => a.id.localeCompare(b.id));
    return devices;
  }
}
