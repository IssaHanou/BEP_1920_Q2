import { Buttons } from "./buttons";

describe("ManageComponent", () => {
  let buttons: Buttons;

  beforeEach(() => {
    buttons = new Buttons();
  });

  it("should create", () => {
    expect(buttons).toBeTruthy();
  });

  it("should be able to add new buttons", () => {
    const jsonData = JSON.parse(`{
          "id": "stop",
          "disabled": true
        }
    `);
    buttons.setButton(jsonData);
    expect(buttons.all.get("stop")).toBeTruthy();
  });

  it("should be able to update existing buttons", () => {
    const jsonData = JSON.parse(`[{
          "id": "stop",
          "disabled": true
        }]
    `);
    buttons.setButtons(jsonData);
    expect(buttons.all.get("stop").getValue()).toBe(true);
    const jsonData2 = JSON.parse(`[{
          "id": "stop",
          "disabled": false
        }]
    `);
    buttons.setButtons(jsonData2);
    expect(buttons.all.get("stop").getValue()).toBe(false);
  });
});
