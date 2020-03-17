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
    this.allSlides = ["Opening",
      "Promenade 1",
      "Gnomus (deel 1)",
      "Promenade 2",
      "Het oude kasteel (deel 2)",
      "Promenade 3",
      "Tuilerieen (deel 3)",
      "Bydlo (deel 4)",
      "Promenade 4",
      "Ballet van de kuikens (deel 5)",
      "Samuel en Schmuyle (deel 6)",
      "Poort van Kiev (deel 10)",
      "Einde"];
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
        if (this.current < this.allSlides.length - 1) {
          this.current++;
        }
        break;
      }
      case "prev": {
        if (this.current > 0) {
          this.current--;
        }
        break;
      }
      case "last": {
        this.current = this.allSlides.length - 1;
        break;
      }
      case "first": {
        this.current = 0;
        break;
      }
      case "specific": {
        if (action.value >= 0 && action.value < this.allSlides.length) {
          this.current = action.value;
        }
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
