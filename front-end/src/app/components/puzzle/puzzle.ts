export class Puzzle {
  id: string;
  status: boolean;
  description: string;

  constructor(id, description) {
    this.id = id;
    this.status = false;
    this.description = description;
  }

  /**
   * Updates status of puzzle.
   */
  public updateStatus(newStatus) {
    this.status = newStatus;
  }
}
