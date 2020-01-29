import { Component, OnInit } from "@angular/core";
import { AppComponent } from "../../app.component";
import { EventsMethods } from "../EventsMethods";

/**
 * The puzzle component controls the puzzles tables and is shown in the "Puzzels" box on the home page.
 */
@Component({
  selector: "app-puzzle",
  templateUrl: "./puzzle.component.html",
  styleUrls: ["./puzzle.component.css", "../../../assets/css/main.css", "./../events.css"],
  providers: [EventsMethods]
})
export class PuzzleComponent implements OnInit {

  /**
   * Hint to send for a certain puzzle.
   */
  hint: string;

  /**
   * Id of device to which a hint must be sent.
   */
  device: string;

  /**
   * Events methods to access methods for events and puzzles.
   */
  methods;

  constructor(private app: AppComponent) {
    this.methods = new EventsMethods(app);
  }

  ngOnInit() {
  }


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
   * When a hint has been chosen for a puzzle and the accompanying "Stuur" button is clicked,
   * the selected hint is sent as instruction to hint devices.
   */
  onSelectedHint() {
    if (
      this.hint !== undefined &&
      this.hint !== ""
    ) {
      this.app.sendInstruction([
        {instruction: "hint", value: this.hint}

      ]);
      this.hint = "";
    }
  }

  /**
   * Get list of devices which can show hints.
   */
  getHintDevices() {
    return [];
  }
}
