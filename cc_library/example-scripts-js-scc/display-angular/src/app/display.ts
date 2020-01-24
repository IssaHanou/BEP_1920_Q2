const Device = require("sciler"); // production

export class Display extends Device {
  private hint: string;
  constructor(config) {
    super(config, function(date, level, message) {
      const formatDate = function(date) {
        return (
          date.getDate() +
          "-" +
          date.getMonth() +
          1 +
          "-" +
          date.getFullYear() +
          " " +
          date.getHours() +
          ":" +
          date.getMinutes() +
          ":" +
          date.getSeconds()
        );
      };
      console.log(
        "time=" + formatDate(date) + " level=" + level + " msg=" + message
      ); // call own logger
    });
    this.hint = "";
  }

  getStatus() {
    return {
      hint: this.hint
    };
  }

  performInstruction(action) {
    switch (action.instruction) {
      case "hint": {
        this.hint = action.value;
        this.statusChanged();
        break;
      }
      default: {
        return false;
      }
    }
    return true;
  }

  test() {
    this.hint = "test";
  }

  reset() {
    this.hint = "";
    this.statusChanged();
  }
}
