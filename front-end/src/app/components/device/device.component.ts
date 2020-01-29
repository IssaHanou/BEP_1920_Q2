import { Component, OnInit, ViewChild } from "@angular/core";
import { Device } from "./device";
import { AppComponent } from "../../app.component";
import { MatSort, MatTableDataSource } from "@angular/material";

/**
 * The device component controls the device table in the "Apparaten" box on the home page.
 */
@Component({
  selector: "app-device",
  templateUrl: "./device.component.html",
  styleUrls: ["./device.component.css", "../../../assets/css/main.css"]
})
export class DeviceComponent implements OnInit {
  
  /**
   * The keys used by the table to retrieve data from the DataSource
   */
  deviceColumns: string[] = ["id", "connection", "component", "status", "test"];

  /**
   * Map keeping track of which rows are collapsed.
   */
  collapsed: Map<string, boolean>;

  /**
   * Control the sorting of the table.
   */
  @ViewChild("DeviceTableSort", { static: true }) sort: MatSort;

  constructor(private app: AppComponent) {
    this.collapsed = new Map<string, boolean>();
  }

  ngOnInit() {}

  /**
   * Returns list of Device objects with their current status and connection.
   * Return in the form of map table data source, with sorting enabled.
   */
  public getDeviceList(): MatTableDataSource<Device> {
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
  getComponents(status: Map<string, any>, deviceId: string): any {
    if (!this.collapsed.get(deviceId)) {
      return "open onderdelen en status";
    } else if (status.size === 0) {
      return "geen status";
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
   * Clicking a row should unfold the components and their statuses.
   * Clicking again should collapse.
   * @param row the row to collapse with full device data
   */
  collapseComponents(row) {
    const oldValue = this.collapsed.get(row.id);
    this.collapsed.set(row.id, !oldValue);
  }

  /**
   * Returns the description of a device with id.
   */
  public getDeviceDescription(id: string) {
    return this.app.deviceList.all.get(id).description;
  }
}
