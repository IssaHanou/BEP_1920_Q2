import { Component, OnInit, ViewChild } from "@angular/core";
import { Device } from "./device";
import { AppComponent } from "../../app.component";
import {MatPaginator, MatSort, MatTableDataSource} from "@angular/material";

@Component({
  selector: "app-device",
  templateUrl: "./device.component.html",
  styleUrls: ["./device.component.css", "../../../assets/css/main.css"]
})
export class DeviceComponent implements OnInit {
  deviceColumns: string[] = ["id", "connection", "component", "status"];

  @ViewChild("DeviceTableSort", {static: true}) sort: MatSort;
  @ViewChild("DevicePaginator", {static: true}) paginator: MatPaginator;

  constructor(private app: AppComponent) {}

  ngOnInit() {
  }

  /**
   * Returns list of Device object with their current status and connection.
   * Return in the form of map table data source, with sorting enabled.
   */
  public getDeviceStatus(): MatTableDataSource<Device> {
    const devices: Device[] = [];
    for (const device of this.app.deviceList.all.values()) {
      devices.push(device);
    }
    devices.sort((a: Device, b: Device) => a.id.localeCompare(b.id));

    const dataSource = new MatTableDataSource<Device>(devices);
    dataSource.sort = this.sort;
    setTimeout(() => dataSource.paginator = this.paginator);
    return dataSource;
  }

  /**
   * Creates list of components (keys of status maps), in alphabetical order.
   * Returns string with each component's value on new line.
   */
  formatStatus(status: Map<string, any>) {
    const keys = Array.from(status.keys());
    keys.sort();

    let result = "";
    keys.forEach((key: string) => {
      const value = status.get(key);
      if (Array.isArray(value)) {
        result += "[";
        for (let i = 0; i < value.length; i++) {
          result += value[i];
          if (i < value.length - 1) {result += ","}
        }
        result += "]\n";
      } else {
        result += value + "\n";
      }
    });
    return result;
  }

  /**
   * Creates list of components (keys of status maps), in alphabetical order.
   * Returns string with each component on new line.
   */
  getComponents(status: Map<string, any>) {
    const keys = Array.from(status.keys());
    keys.sort();
    let result = "";
    keys.forEach((key: string) => {result += key + "\n"});
    return result;
  }
}
