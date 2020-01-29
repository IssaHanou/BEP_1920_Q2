import { Component } from "@angular/core";
import { AppComponent } from "../../app.component";
import { Hint } from "./hint";

/**
 * The hint component controls the sending of hints in the "Hint" box on the home pgae.
 */
@Component({
  selector: "app-hint",
  templateUrl: "./hint.component.html",
  styleUrls: ["./hint.component.css", "../../../assets/css/main.css"]
})
export class HintComponent {
  /**
   * The contents of the hint text field.
   */
  customHint: string;

  /**
   * The contents of the hint selection box.
   */
  predefinedHint: string;

  constructor(private app: AppComponent) {}

  /**
   * List of puzzle names to use as dividers in the selection box for predefined hints.
   */
  getPuzzleList(): Hint[] {
    const list = [];
    for (const hint of this.app.hintList) {
      list.push(hint);
    }
    return list;
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
   * When predefined hint has been chosen and the accompanying "Stuur" button is clicked,
   * the selected hint is sent as instruction to hint devices.
   */
  onPredefinedHint() {
    if (
      this.predefinedHint !== undefined &&
      this.predefinedHint !== "" &&
      this.predefinedHint !== "---"
    ) {
      this.app.sendInstruction([
        {
          instruction: "hint",
          value: this.predefinedHint,
          topic: "hint"
        }
      ]);
      this.predefinedHint = "---";
    }
  }

  /**
   * When custom hint has been typed and the accompanying "Stuur" button is clicked,
   * the typed hint is sent as instruction to hint devices.
   */
  onCustomHint() {
    if (this.customHint !== undefined && this.customHint !== "") {
      this.app.sendInstruction([
        {
          instruction: "hint",
          value: this.customHint,
          topic: "hint"
        }
      ]);
      this.customHint = "";
    }
  }
}
