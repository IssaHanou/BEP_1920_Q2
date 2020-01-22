import { Timer } from "./timer";

/**
 * Class keeping track of the timers, through use of map with timer id's.
 */
export class Timers {
  all: Map<string, Timer>;

  constructor() {
    this.all = new Map<string, Timer>();
  }

  /**
   * Set the new duration and state of a certain timer.
   * If the timer did not yet exist, create a new one.
   */
  setTimer(jsonData) {
    if (this.all.has(jsonData.id)) {
      this.all.get(jsonData.id).update(jsonData.duration, jsonData.state);
    } else {
      this.all.set(jsonData.id, new Timer(jsonData));
    }
  }

  /**
   * Return timer with id t.
   */
  getTimer(t) {
    if (this.all.has(t)) {
      return this.all.get(t);
    } else {
      return null;
    }
  }

  /**
   * Return the map of timers.
   */
  getAll() {
    return this.all;
  }
}
