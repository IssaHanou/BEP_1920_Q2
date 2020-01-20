export class Button {

  id: string;
  disabled: boolean;

  constructor(jsonData) {
    this.id = jsonData.id;
    this.disabled = jsonData.disabled;
  }

  /**
   * Update with new disabled value.
   */
  updateDisabled(newValue) {
    this.disabled = newValue;
  }

  /**
   * Return the disabled value of the button.
   */
  getValue() {
    return this.disabled;
  }
}
