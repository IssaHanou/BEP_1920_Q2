export class Puzzle {
  id: string;
  status: Map<string, boolean>;
  description: string;

  constructor(jsonData) {
    this.id = jsonData["id"];
    this.description = jsonData["description"]
    // TODO currently shows status per rule not puzzle
    this.status = jsonData["status"];
  }

  /**
   * Updates status of puzzle.
   */
  updateStatus(newStatus) {
    this.status = newStatus;
  }
}
