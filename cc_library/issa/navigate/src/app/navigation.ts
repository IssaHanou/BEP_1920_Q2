const Device = require("sciler"); // production

export class Navigation extends Device {
  private current: number;
  private slides: string[];

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
    this.current = 0;
    this.slides = ["Opening", "Promenade 1", "Deel 1", "Promenade 2", "Deel 2", "Promenade 3", "Deel 3", "Einde"];
  }

  getStatus() {
    return {
      showing: this.current
    };
  }

  performInstruction(action) {
    switch (action.instruction) {
      case "set": {
        this.current = action.value;
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
    this.current = this.slides.length - 1;
  }

  reset() {
    this.current = 0;
  }
}
