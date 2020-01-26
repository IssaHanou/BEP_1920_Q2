import { Button } from "./button";

describe("ManageComponent", () => {
  let button: Button;

  beforeEach(() => {
    const jsonData = JSON.parse(`{
          "id": "stop",
          "disabled": true
        }
    `);
    button = new Button(jsonData);
  });

  it("should create", () => {
    expect(button).toBeTruthy();
  });

  it("should set value", () => {
    expect(button.disabled).toBe(true);
    button.updateDisabled(false);
    expect(button.disabled).toBe(false);
  });
});
