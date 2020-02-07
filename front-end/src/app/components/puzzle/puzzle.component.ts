import { Component, OnInit } from "@angular/core";
import { AppComponent } from "../../app.component";
import { EventsMethods } from "../EventsMethods";

/**
 * The puzzle component controls the puzzles tables and is shown in the "Puzzels" box on the home page.
 */
@Component({
  selector: "app-puzzle",
  templateUrl: "./puzzle.component.html",
  styleUrls: [
    "./puzzle.component.css",
    "../../../assets/css/main.css",
    "./../events.css"
  ],
  providers: [EventsMethods]
})
export class PuzzleComponent implements OnInit {
  /**
   * Hint to send for a certain puzzle.
   */
  hint: string;

  /**
   * Topic to which a hint must be sent.
   */
  topic: string;

  constructor(private app: AppComponent, private methods: EventsMethods) {
    this.hint = "";
    this.topic = "";
  }

  ngOnInit() {}

  /**
   * Hint list used for selection of predefined hints.
   * This is generated each time from the app hint list, to ensure updated version.
   */
  getHintList(puzzle: string): string[] {
    const list = [];
    for (const obj of this.app.hintList) {
      if (obj.puzzle === puzzle) {
        for (const hint of obj.hints) {
          list.push(hint);
        }
        return list;
      }
    }
  }

  /**
   * When a hint has been chosen for a puzzle, a device (topic) has been chosen to send to
   * and the accompanying send button is clicked,
   * the selected hint is sent as instruction to hint devices.
   */
  onSelectedHint(puzzleName: string) {
    if (this.hint !== undefined && this.hint !== "") {
      this.app.sendHint(this.hint, this.topic, puzzleName);
    }
    this.hint = "";
    this.topic = "";
  }
}
