import { Device } from "./device";

/**
 * Class keeping track of the devices, through use of map with device id's.
 */
export class Devices {
  /**
   * Maps device id to device object.
   */
  all: Map<string, Device>;
  /**
   * Maps label to all devices that listen to that label.
   */
  labels: Map<string, string[]>;

  constructor() {
    this.all = new Map<string, Device>();
    this.labels = new Map<string, string[]>();
  }

  /**
   * Create devices initially from json data with id, description and labels.
   */
  createDevices(jsonData) {
    for (const device of jsonData) {
      this.all.set(device.id, new Device(device));
    }
  }

  /**
   * Create the label map with for each label a list of devices that listen to it.
   */
  createLabelMap() {
    for (const deviceId of this.all.keys()) {
      const device = this.all.get(deviceId);
      for (const index in device.labels) {
        const label = device.labels[0];
        if (this.labels.has(label)) {
          this.labels.get(label).push(deviceId);
        } else {
          this.labels.set(label, [deviceId]);
        }
      }
    }
  }

  /**
   * Set the new status and connection state of a certain device.
   * If the device did not yet exist, create a new one.
   */
  setDevice(jsonData) {
    if (this.all.has(jsonData.id)) {
      this.all.get(jsonData.id).updateStatus(jsonData.status);
      this.all.get(jsonData.id).updateConnection(jsonData.connection);
    } else {
      this.all.set(jsonData.id, new Device(jsonData));
    }
  }

  /**
   * Update the status of device front-end.
   * Update the component with id to status.
   */
  updateDevice(id, status) {
    if (this.all.has("front-end")) {
      const newStatus = {};
      newStatus[id] = status;
      this.all.get("front-end").updateStatus(newStatus);
    }
  }

  /**
   * Return device with id dev.
   */
  getDevice(dev: string) {
    if (this.all.has(dev)) {
      return this.all.get(dev);
    } else {
      return null;
    }
  }
}
