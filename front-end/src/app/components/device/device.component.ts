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
  deviceColumns: string[] = ["id", "connection", "unfold", "component", "status", "test"];

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
        this.collapsed.set(device.id, true);
      }
    }
    devices.sort((a: Device, b: Device) => a.id.localeCompare(b.id));

    const dataSource = new MatTableDataSource<Device>(devices);
    dataSource.sort = this.sort;
    return dataSource;
  }

  /**
   * In the table, show the front-end as device with name "operator scherm".
   */
  getName(deviceId: string): string {
    if (deviceId === "front-end") {
      return "operator scherm";
    } else {
      return deviceId;
    }
  }

  /**
   * Creates list of components (keys of status maps), in alphabetical order.
   * Returns string with each component on new line.
   * If components are collapsed, return default.
   * For the front-end, only show gameState.
   */
  getComponents(status: Map<string, any>, deviceId: string): any {
    if (this.collapsed.get(deviceId)) {
      return "";
    } else if (deviceId === "front-end") {
      return "gameState";
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
   * For the front-end only show gameState status.
   */
  formatStatus(status: Map<string, any>, deviceId: string): string {
    if (this.collapsed.get(deviceId)) {
      return "";
    } else if (deviceId === "front-end") {
      return status.get("gameState");
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
   * The test button can not be clicked when game is in play (timer is running).
   */
  testDevice(deviceId: string) {
    this.app.sendInstruction([
      { instruction: "test device", device: deviceId }
    ]);
  }

  /**
   * Clicking the 'status' button should show the components and their statuses.
   * Clicking the 'sluit' button should hide the components and their statuses.
   * @param deviceID the status of this device should be shown/hidden
   */
  collapseComponents(deviceID: string) {
    const oldValue = this.collapsed.get(deviceID);
    this.collapsed.set(deviceID, !oldValue);
  }
}
