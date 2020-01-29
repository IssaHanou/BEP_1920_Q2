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
  hint: string = "";

  /**
   * Id of device to which a hint must be sent.
   */
  device: string = "";

  constructor(private app: AppComponent, private methods: EventsMethods) {}

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
   * When a hint has been chosen for a puzzle and the accompanying "Stuur" button is clicked,
   * the selected hint is sent as instruction to hint devices.
   */
  onSelectedHint() {
    if (this.hint !== undefined && this.hint !== "") {
      if (this.device === "all devices" || this.device === "") {
        this.device = "hint"; // hints to all devices should be published to topic hint
      }
      this.app.sendInstruction([
        {
          instruction: "hint",
          value: this.hint,
          topic: this.device
        }
      ]);
      this.hint = "";
      this.device = "";
    }
  }

  /**
   * Get list of devices which can show hints, and add the option to send the hint to `all devices`.
   */
  getHintDevices(): string[] {
    let hintDevices: string[] = ["all devices"];
    if (this.app.deviceList.labels.has("hint")) {
      hintDevices = hintDevices.concat(this.app.deviceList.labels.get("hint"));
    }
    return hintDevices;
  }
}
