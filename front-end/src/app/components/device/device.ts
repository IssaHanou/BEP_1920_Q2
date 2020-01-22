/**
 * Puzzle object, which has an id, status (can be anything, depending on device implementation) and connection status (boolean).
 */
export class Device {
  id: string;
  status: Map<string, any>;
  connection: boolean;

  constructor(jsonData) {
    this.id = jsonData.id;
    this.updateConnection(jsonData.connection);
    this.status = new Map<string, any>();
    this.updateStatus(jsonData.status);
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
        this.status.delete(key);
        this.status.set(key, jsonStatus[key]);
      } else {
        this.status.set(key, jsonStatus[key]);
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
