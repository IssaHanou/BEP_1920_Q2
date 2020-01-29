import { Event } from "./event";

/**
 * Class keeping track of the events, through use of map with puzzle id's.
 */
export class Events {
  /**
   * Map with rule id keys to event values.
   */
  rules: Map<string, Event>;
  /**
   * Map with event name keys to event lis values.
   */
  rulesPerEvent: Map<string, Event[]>;

  constructor() {
    this.rules = new Map<string, Event>();
    this.rulesPerEvent = new Map<string, Event[]>();
  }

  /**
   * Set the new status of a certain event.
   * If the event did not yet exist, create a new one with its description.
   */
  updatePuzzles(jsonData) {
    for (const object of jsonData) {
      if (!this.rules.has(object.id)) {
        this.rules.set(object.id, new Event(object.id, object.description, object.status, object.eventName, object.puzzle));
      }
      this.rules.get(object.id).updateStatus(object.status);
    }
  }

  /**
   * Creates a map where the name of a puzzle/general event is mapped
   * to a list of event objects that belong to that event/puzzle.
   */
  createRulesPerEvent() {
    for (const ruleId of this.rules.keys()) {
      const rule = this.rules.get(ruleId);
      if (this.rulesPerEvent.has(rule.puzzleName)) {
        this.rulesPerEvent.get(rule.puzzleName).push(rule);
      } else {
        this.rulesPerEvent.set(rule.puzzleName, [rule]);
      }
    }
  }
}
