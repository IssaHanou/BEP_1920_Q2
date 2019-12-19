import { Timer } from "./timer";

export class Timers {
  all: Map<string, Timer>;

  constructor() {
    this.all = new Map<string, Timer>();
  }

  setTimer(jsonData) {
    if (this.all.has(jsonData.id)) {
      this.all.get(jsonData.id).update(jsonData.status, jsonData.state);
    } else {
      this.all.set(jsonData.id, new Timer(jsonData));
    }
  }

  getTimer(t) {
    if (this.all.has(t)) {
      return this.all.get(t);
    } else {
      return null;
    }
  }

  getAll() {
    return this.all;
  }
}
