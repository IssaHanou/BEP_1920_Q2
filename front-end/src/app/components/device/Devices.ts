import {Device} from "./Device";

/**
 * Devices has a Map all containing all devices with a key that is the same as the id.
 */
export class Devices {
  all: Map<string, Device>;

  constructor() {
    this.all = new Map<string,Device>()
  }

  /**
   * setDevice either updates an existing Device with the update methods or creates a new one.
   * @param jsonData json object with keys id, status and connection.
   */
  setDevice(jsonData){
    if (this.all.has(jsonData.id)) {
      this.all.get(jsonData.id).updateStatus(jsonData.status);
      this.all.get(jsonData.id).updateConnection(jsonData.connection)
    } else {
      this.all.set(jsonData.id, new Device(jsonData))
    }

  }
}
