import { MatTableDataSource } from "@angular/material";
import {Event} from "./event/event";
import {AppComponent} from "../app.component";
import {Inject} from "@angular/core";

export class EventsMethods {

  /**
   * The keys used by the table to retrieve data from the DataSource
   */
  ruleColumns: string[] = ["id", "status", "done"];

  /**
   * columns for the puzzle table.
   */
  puzzleColumns: string[] = ["puzzle"];

  constructor(@Inject(AppComponent) private app: AppComponent) {}

  /**
   * Returns list of event names.
   * Return in the form of map table data source, with sorting enabled.
   */
  public getEvents(shouldBePuzzle: boolean): MatTableDataSource<string> {
    const events: string[] = [];
    for (const eventName of this.app.puzzleList.rulesPerEvent.keys()) {
      if (this.app.puzzleList.rulesPerEvent.get(eventName).length > 0 &&
        this.app.puzzleList.rulesPerEvent.get(eventName)[0].isPuzzle == shouldBePuzzle) {
        events.push(eventName);
      }
    }
    events.sort((a, b) => a.localeCompare(b));
    return new MatTableDataSource<string>(events);
  }

  /**
   * Returns list of Event objects that belong to the puzzle name with their current status.
   * Return in the form of map table data source, with sorting enabled.
   */
  public getRulesPerEvent(puzzleName: string): MatTableDataSource<Event> {
    const puzzles: Event[] = [];
    for (const rule of this.app.puzzleList.rulesPerEvent.get(puzzleName)) {
      puzzles.push(rule);
    }
    puzzles.sort((a: Event, b: Event) => a.id.localeCompare(b.id));
    return  new MatTableDataSource<Event>(puzzles);
  }


  /**
   * When button in the events or puzzles table is pressed, manually override the finished status of rule in back-end.
   */
  public finishRule(ruleId: string) {
    this.app.sendInstruction([{ instruction: "finish rule", rule: ruleId }]);
  }

  /**
   * Returns the description of an event with id eventID.
   */
  public getEventDescription(eventID: string) {
    return this.app.puzzleList.rules.get(eventID).description;
  }
}
