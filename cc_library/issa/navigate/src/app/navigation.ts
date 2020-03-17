const Device = require("sciler"); // production

export class Navigation extends Device {
  private current: number;
  private slides: string[];

  private next: boolean;
  private prev: boolean;
  private first: boolean;
  private last: boolean;
  private specific: number;
  private specificSubmitted: boolean;

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
    this.slides = ["Opening",
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

    this.prev = false;
    this.next = false;
    this.last = false;
    this.first = false;
    this.specific = 0;
    this.specificSubmitted = false;
  }

  getStatus() {
    return {
      showing: this.current,
      prev: this.prev,
      next: this.next,
      first: this.first,
      last: this.last,
      specific: this.specific,
      specificSubmitted: this.specificSubmitted
    };
  }

  performInstruction(action) {
    return false;
  }

  test() {
    this.current = this.slides.length - 1;
  }

  reset() {
    this.current = 0;
  }

  sendInstruction(instruction: string, value: number) {
    console.log("sending instruction: " + instruction + " " + value);
    switch (instruction) {
      case "first":
        this.first = true;
        this.statusChanged();
        this.first = false;
        this.statusChanged();
        break;
      case "last":
        this.last = true;
        this.statusChanged();
        this.last = false;
        this.statusChanged();
        break;
      case "next":
        this.next = true;
        this.statusChanged();
        this.next = false;
        this.statusChanged();
        break;
      case "prev":
        this.prev = true;
        this.statusChanged();
        this.prev = false;
        this.statusChanged();
        break;
      case "specific":
        this.specific = value;
        this.specificSubmitted = true;
        this.statusChanged();
        this.specificSubmitted = false;
        this.statusChanged();
        break;
      default:
        this.log("warning", "illegal instruction: " + instruction);
        break;
    }
  }
}
