import { Component, OnInit, ViewChild } from "@angular/core";
import { Device } from "./device";
import { AppComponent } from "../../app.component";
import { MatSort, MatTableDataSource } from "@angular/material";

@Component({
  selector: "app-device",
  templateUrl: "./device.component.html",
  styleUrls: ["./device.component.css", "../../../assets/css/main.css"]
})
export class DeviceComponent implements OnInit {
  deviceColumns: string[] = ["id", "connection", "component", "status", "test"];
  collapsed: Map<string, boolean>;

  @ViewChild("DeviceTableSort", { static: true }) sort: MatSort;

  constructor(private app: AppComponent) {
    this.collapsed = new Map<string, boolean>();
  }

  ngOnInit() {}

  /**
   * Returns list of Device objects with their current status and connection.
   * Return in the form of map table data source, with sorting enabled.
   */
  public getDeviceStatus(): MatTableDataSource<Device> {
    const devices: Device[] = [];
    for (const device of this.app.deviceList.all.values()) {
      devices.push(device);
      if (!this.collapsed.has(device.id)) {
        this.collapsed.set(device.id, false);
      }
    }
    devices.sort((a: Device, b: Device) => a.id.localeCompare(b.id));

    const dataSource = new MatTableDataSource<Device>(devices);
    dataSource.sort = this.sort;
    return dataSource;
  }

  /**
   * Creates list of components (keys of status maps), in alphabetical order.
   * Returns string with each component on new line.
   * If components are collapsed, return default.
   */
  getComponents(status: Map<string, any>, deviceId: string): string {
    if (!this.collapsed.get(deviceId)) {
      return "click to see components";
    } else if (status.size == 0) {
      return "nothing to show";
    } else {
      const keys = Array.from(status.keys());

      keys.sort();
      let result = "";
      keys.forEach((key: string) => {
        result += key + "\n";
      });
      return result;
    }
  }

  /**
   * Creates list of components (keys of status maps), in alphabetical order.
   * Returns string with each component's value on new line.
   * If components are collapse, return nothing.
   */
  formatStatus(status: Map<string, any>, deviceId: string): string {
    if (!this.collapsed.get(deviceId)) {
      return "";
    } else {
      const keys = Array.from(status.keys());
      keys.sort();

      let result = "";
      keys.forEach((key: string) => {
        const value = status.get(key);
        if (Array.isArray(value)) {
          result += "[";
          for (let i = 0; i < value.length; i++) {
            result += value[i];
            if (i < value.length - 1) {
              result += ",";
            }
          }
          result += "]\n";
        } else {
          result += value + "\n";
        }
      });
      return result;
    }
  }

  /**
   * When button is pressed, test a single device.
   */
  testDevice(deviceId: string) {
    this.app.sendInstruction([
      { instruction: "test device", device: deviceId }
    ]);
  }

  /**
   * Clicking a row should collapse the components and their statuses.
   * Clicking again should unfold.
   * @param row the row to collapse with full device data
   */
  collapseComponents(row) {
    // document.getElementById(row.id + "-status").style;
    const oldValue = this.collapsed.get(row.id);
    this.collapsed.set(row.id, !oldValue);
  }
}
