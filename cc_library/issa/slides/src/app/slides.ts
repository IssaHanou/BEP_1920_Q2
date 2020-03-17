const Device = require("sciler"); // production

export class Slides extends Device {
  private current: number;
  private allSlides: string[];

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
    this.allSlides = ["Opening", "Promenade 1", "Deel 1", "Promenade 2", "Deel 2", "Promenade 3", "Deel 3", "Einde"];
    this.current = 0;
  }

  getStatus() {
    return {
      current: this.current
    };
  }

  performInstruction(action) {
    switch (action.instruction) {
      case "next": {
        this.current++;
        break;
      }
      case "prev": {
        this.current--;
        break;
      }
      case "last": {
        this.current = this.allSlides.length - 1;
        break;
      }
      case "specific": {
        this.current = action.value;
        break;
      }
      default: {
        return false;
      }
    }
    return true;
  }

  test() {
    this.current = this.allSlides.length - 1;
  }

  reset() {
    this.current = 0;
    this.statusChanged();
  }
}
