import { Component, OnInit, ViewChild } from "@angular/core";
import { Comp, Device } from "./device";
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
   * The keys used by the device table to retrieve data from the DataSource.
   */
  deviceColumns: string[] = ["id", "connection", "unfold", "test"];
  /**
   * The keys used by the component table to retrieve data from the component list.
   */
  componentColumns: string[] = ["component", "status"];

  /**
   * Control the sorting of the table.
   */
  @ViewChild("DeviceTableSort", { static: true }) sort: MatSort;

  constructor(private app: AppComponent) {}

  ngOnInit() {}

  /**
   * Returns list of Device objects with their current status and connection.
   * Return in the form of map table data source, with sorting enabled.
   */
  public getDeviceList(): MatTableDataSource<Device> {
    const devices: Device[] = [];
    for (const device of this.app.deviceList.all.values()) {
      devices.push(device);
    }
    devices.sort((a: Device, b: Device) => a.id.localeCompare(b.id));

    const dataSource = new MatTableDataSource<Device>(devices);
    dataSource.sort = this.sort;
    return dataSource;
  }

  /**
   * Returns list of Comp objects with their current status.
   * For the front-end, only return the game state not the button pressed statuses.
   */
  public getComponentList(deviceId: string): Comp[] {
    const status = this.app.deviceList.getDevice(deviceId).status;
    if (deviceId === "front-end") {
      return [status.get("gameState")];
    }
    const ret = [];
    const keys = Array.from(status.keys());
    keys.sort();
    keys.forEach((key: string) => {
      ret.push(status.get(key));
    });
    return ret;
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
   * When button is pressed, test a single device.
   * The test button can not be clicked when game is in play (timer is running).
   */
  testDevice(deviceId: string) {
    this.app.sendInstruction([
      { instruction: "test device", device: deviceId }
    ]);
  }
}
