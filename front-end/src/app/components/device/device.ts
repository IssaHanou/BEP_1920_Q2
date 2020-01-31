/**
 * Event object, which has an id, status (can be anything, depending on device implementation) and connection status (boolean).
 */
export class Device {
  id: string;
  status: Map<string, Comp>;
  connection: boolean;
  description: string;
  labels: string[];

  constructor(jsonData) {
    this.id = jsonData.id;
    this.labels = jsonData.labels;
    this.description = jsonData.description;
    this.status = new Map<string, Comp>();
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
    const keys = Object.keys(jsonStatus);
    for (const key of keys) {
      if (this.status.has(key)) {
        this.status.get(key).status = jsonStatus[key];
      } else {
        this.status.set(key, new Comp(key, jsonStatus[key]));
      }
    }
  }

  /**
   * getValue returns status of specific component
   * @param comp component id
   */
  getValue(comp) {
    return this.status.get(comp);
  }
}

/**
 * Comp class contains information about a component: its id and status.
 */
export class Comp {
  id: string;
  status: any;

  constructor(id, status) {
    this.id = id;
    this.status = status;
  }
}
