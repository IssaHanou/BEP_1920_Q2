import {Device} from "./device";
import {DeviceComponent} from "./device.component";


describe("DeviceComponent", () => {
  let device: Device;

  beforeEach(() => {
    const jsonData = JSON.parse(`{
          "id": "Door",
          "status": {"door": true},
           "connection": false
        }
    `);
    device = new Device(jsonData);
  });

  it("should create", () => {
    expect(device).toBeTruthy();
  });

  it("should set connection", () => {
    expect(device.connection).toBe(false);
    device.updateConnection(true);
    expect(device.connection).toBe(true);
    }
  );

  it("should set status", () => {
    console.log(device)
      expect(device.getValue("door")).toBe(true);
    const jsonData = JSON.parse(`{ "door" : false }`);
      device.updateStatus(jsonData);
      expect(device.connection).toBe(false);
    }
  );

});
