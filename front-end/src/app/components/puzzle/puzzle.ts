export class Puzzle {
  id: string;
  status: boolean;
  description: string;
  keys = ["id", "status", "description"];

  constructor(jsonData) {
    this.id = jsonData[this.keys[0]];
    this.description = jsonData[this.keys[2]];
    // TODO currently shows status per rule not puzzle
    this.status = jsonData[this.keys[1]];
  }

  /**
   * Updates status of puzzle.
   */
  updateStatus(newStatus) {
    this.status = newStatus;
  }
}
