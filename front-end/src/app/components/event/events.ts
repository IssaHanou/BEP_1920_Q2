import { Event } from "./event";

/**
 * Class keeping track of the events, through use of map with puzzle id's.
 */
export class Events {
  all: Map<string, Event>;

  constructor() {
    this.all = new Map<string, Event>();
  }

  /**
   * Set the new status of a certain event.
   * If the event did not yet exist, create a new one with its description.
   */
  updatePuzzles(jsonData) {
    for (let i = 0; i < jsonData.length; i++) {
      const object = jsonData[i];
      if (!this.all.has(object.id)) {
        this.all.set(object.id, new Event(object.id, object.description, object.status, object.puzzle));
      }
      this.all.get(object.id).updateStatus(object.status);
    }
  }
}
