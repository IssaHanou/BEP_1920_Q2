import { Puzzle } from "./puzzle";

export class Puzzles {
  all: Map<string, Puzzle>;

  constructor() {
    this.all = new Map<string, Puzzle>();
  }

  /**
   * Receives list of json objects (keys id and status)
   */
  updatePuzzles(jsonData) {
    for (const object of jsonData) {
      if (this.all.has(object.id)) {
        this.all.get(object.id).updateStatus(object.status);
      } else {
        this.all.set(object.id, new Puzzle(object.id, object.description));
      }
    }
  }

  /**
   * Set the puzzles in the list with their id and description
   * @param events map with id keys and description values
   */
  setPuzzles(events: Map<string, string>) {
    for (const rule of events) {
      this.all.set(rule[0], new Puzzle(rule[0], rule[1]));
    }
  }
}
