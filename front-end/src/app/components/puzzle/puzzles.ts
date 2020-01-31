import { Puzzle } from "./puzzle";

/**
 * Class keeping track of the puzzles, through use of map with puzzle id's.
 */
export class Puzzles {
  all: Map<string, Puzzle>;

  constructor() {
    this.all = new Map<string, Puzzle>();
  }

  /**
   * Set the new status of a certain puzzle.
   * If the timer did not yet exist, create a new one with its description.
   */
  updatePuzzles(jsonData) {
    for (const object of jsonData) {
      if (!this.all.has(object.id)) {
        this.all.set(object.id, new Puzzle(object.id, object.description));
      }
      this.all.get(object.id).updateStatus(object.status);
    }
  }

  /**
   * Set the puzzles in the list with their id and description
   * @param id of new puzzle
   * @param description of new puzzle
   */
  addPuzzle(id: string, description: string) {
    this.all.set(id, new Puzzle(id, description));
  }
}
