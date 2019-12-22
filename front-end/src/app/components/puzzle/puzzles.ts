import {Puzzle} from "./puzzle";

export class Puzzles {
  all: Map<string, Puzzle>;

  constructor() {
    this.all = new Map<string, Puzzle>();
  }

  /**
   * Receives list of json objects (keys 'id' and 'status')
   */
  updatePuzzles(jsonData) {
    for (const object of jsonData) {
      if (this.all.has(object["id"])) {
        this.all.get(object["id"]).updateStatus(object["status"]);
      } else {
        this.all.set(object["id"], object);
      }
    }
    console.log(this.all);
  }
}
