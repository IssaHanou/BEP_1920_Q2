/**
 * Event object, which has an id, status (boolean whether it has been solved), description and boolean whether it belongs to puzzle.
 */
export class Event {
  id: string;
  status: boolean;
  description: string;
  isPuzzle: boolean;

  constructor(id, description, status, isPuzzle) {
    this.id = id;
    this.status = status;
    this.description = description;
    this.isPuzzle = isPuzzle;
  }

  /**
   * Updates status of event.
   */
  public updateStatus(newStatus) {
    this.status = newStatus;
  }
}
