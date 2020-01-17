import { Button } from "./button";

/**
 * Devices has a Map all containing all devices with a key that is the same as the id.
 */
export class Buttons {
  all: Map<string, Button>;

  constructor() {
    this.all = new Map<string, Button>();
  }

  /**
   * For every button in the data, update its disabled valued.
   * @param jsonData
   */
  setButtons(jsonData) {
    for (const object of jsonData) {
      this.setButton(object);
    }
  }

  /**
   * setButton either updates an existing Button with the the new disabled status or creates a new button.
   * @param jsonData json object with keys id and disabled.
   */
  setButton(jsonData) {
    if (this.all.has(jsonData.id)) {
      this.all.get(jsonData.id).updateDisabled(jsonData.disabled);
    } else {
      this.all.set(jsonData.id, new Button(jsonData));
    }
  }

  /**
   * getButton is a getter for buttons
   * @param btn button id
   */
  getButton(btn: string) {
    if (this.all.has(btn)) {
      return this.all.get(btn);
    } else {
      return null;
    }
  }
}
