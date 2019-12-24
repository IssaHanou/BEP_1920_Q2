export class Puzzle {
  id: string;
  status: boolean;
  description: string;

  constructor(jsonData) {
    this.id = jsonData.id;
    this.status = jsonData.status;
    this.description = jsonData.description;
  }

  /**
   * Updates status of puzzle.
   */
  public updateStatus(newStatus) {
    this.status = newStatus;
  }
}
