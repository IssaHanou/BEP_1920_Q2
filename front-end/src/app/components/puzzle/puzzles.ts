import {Puzzle} from "./puzzle";

export class Puzzles {
  all: Map<string, Puzzle>;
  keys = ["id", "status", "description"];

  constructor() {
    this.all = new Map<string, Puzzle>();
  }

  /**
   * Receives list of json objects (keys 'id' and 'status')
   */
  updatePuzzles(jsonData) {
    for (const object of jsonData) {
      if (this.all.has(object[this.keys[0]])) {
        this.all.get(object[this.keys[0]]).updateStatus(object[this.keys[1]]);
      } else {
        this.all.set(object[this.keys[0]], object);
      }
    }
    console.log(this.all);
  }
}
