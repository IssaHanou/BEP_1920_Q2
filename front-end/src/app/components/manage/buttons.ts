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
   * @param jsonData of all button object
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
}
