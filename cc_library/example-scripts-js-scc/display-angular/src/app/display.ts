const Device = require("../../../../js-scc"); // development
// const SccLib = require("js-scc"); // production

export class Display extends Device {
  private hint: string;
  private time: string;
  private roomName: string;

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
    this.time = "00:30:00";
    this.roomName = "Escape room";
  }

  getStatus() {
    return {
      "hint": this.hint
    }
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
