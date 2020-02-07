/**
 * Event object, which has an id, status (can be anything, depending on device implementation) and connection status (boolean).
 */
export class Device {
  id: string;
  statusMap: Map<string, Comp>;
  connection: boolean;
  description: string;
  labels: string[];

  constructor(jsonData) {
    this.id = jsonData.id;
    this.labels = jsonData.labels;
    this.description = jsonData.description;
    this.statusMap = new Map<string, Comp>();
    if (jsonData.hasOwnProperty("status")) {
      this.updateStatus(jsonData.status);
    }
    if (jsonData.hasOwnProperty("connection")) {
      this.updateConnection(jsonData.connection);
    }
  }

  /**
   * updateConnection is called on every status update to update the connections status.
   * @param connection boolean
   */
  updateConnection(connection) {
    this.connection = connection;
  }

  /**
   * updateStatus is called on every status update to update the component status.
   * @param jsonStatus json object containing components as key with status as value.
   */
  updateStatus(jsonStatus) {
    for (const key of Object.keys(jsonStatus)) {
      if (this.statusMap.has(key)) {
        this.statusMap.get(key).componentStatus = jsonStatus[key];
      } else {
        this.statusMap.set(key, new Comp(key, jsonStatus[key]));
      }
    }
  }

  /**
   * getValue returns status of specific component
   * @param comp component id
   */
  getValue(comp) {
    return this.statusMap.get(comp);
  }
}

/**
 * Comp class contains information about a component: its id and status.
 */
export class Comp {
  id: string;
  componentStatus: any;

  constructor(id, status) {
    this.id = id;
    this.componentStatus = status;
  }
}
