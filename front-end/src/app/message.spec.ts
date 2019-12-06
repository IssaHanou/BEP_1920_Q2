import { Component, OnInit } from "@angular/core";
import { Message } from "./message";


describe("AppComponent",() => {

  let message: Message;

  beforeEach(() => {
    message = new Message("front-end",
      "confirmation",
      new Date(2019, 1, 5, 14, 16, 1),
      {
        "completed": true,
        "instructed": {
          "device_id": "door",
          "time_sent": "10-05-2019 15:09:14",
          "type": "instruction",
          "contents": {"instruction": "start"}
        }
      })
  });

  it("should create", () => {
    expect(message).toBeTruthy();
  });

  it("should format device's time sent correctly", () => {
    expect(message.timeSent).toEqual("05-01-2019 14:16:01");
  });

  it("should deserialize json correctly", () => {
    let jsonMsg = "{'device_id': 'front-end', " +
      "'type': 'confirmation', " +
      "'time_sent': '05-01-2019 14:16:01', " +
      "'contents': {'completed': true, 'instructed': " +
      "{'device_id': 'door', 'type': 'instruction', 'time_sent': '10-05-2019 15:09:14', " +
      "'contents': {'instruction': 'start'}}}}";
    jsonMsg = replaceAll(jsonMsg,"'", "\"");
    console.log(jsonMsg);
    expect(Message.deserialize(jsonMsg)).toEqual(message);
  });

  function replaceAll(str, find, replace) {
    return str.replace(new RegExp(find, 'g'), replace);
  }
});
