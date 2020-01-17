export class Button {

  id: string;
  disabled: boolean;

  constructor(jsonData) {
    this.id = jsonData.id;
    this.disabled = jsonData.disabled;
  }

  updateDisabled(newValue) {
    this.disabled = newValue;
  }
}
