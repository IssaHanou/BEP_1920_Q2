import { Devices } from "./devices";

describe("DeviceComponent", () => {
  let devices: Devices;

  beforeEach(() => {
    devices = new Devices();
  });

  it("should create", () => {
    expect(devices).toBeTruthy();
  });

  it("should be able to add new devices", () => {
    expect(devices.getDevice("Door")).toBeNull();
    const jsonData = JSON.parse(`{
          "id": "Door",
          "status": {"door": true},
           "connection": false
        }
    `);
    devices.setDevice(jsonData);
    expect(devices.getDevice("Door")).toBeTruthy();
  });

  it("should be able to update existing devices", () => {
    expect(devices.getDevice("Door")).toBeNull();
    const jsonData = JSON.parse(`{
          "id": "Door",
          "status": {"door": true},
           "connection": false
        }
    `);
    devices.setDevice(jsonData);
    expect(devices.getDevice("Door").getValue("door").componentStatus).toBe(true);
    const jsonData2 = JSON.parse(`{
          "id": "Door",
          "status": {"door": false},
           "connection": false
        }
    `);
    devices.setDevice(jsonData2);
    expect(devices.getDevice("Door").getValue("door").componentStatus).toBe(false);
  });
});
