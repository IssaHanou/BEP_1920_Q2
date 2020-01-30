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
  hint: string;

  /**
   * Topic to send hint to.
   */
  topic: string;

  constructor(private app: AppComponent) {}

  /**
   * List of hints that have been sent.
   */
  getHintLog(): string[] {
    const list = [];
    for (const hint of this.app.sentHints) {
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
   * When a hint has been entered, a device (topic) has been chosen to send to
   * and the accompanying send button is clicked,
   * the selected hint is sent as instruction to hint devices.
   */
  onCustomHint() {
    if (this.hint !== undefined && this.hint !== "") {
      this.app.sendHint(this.hint, this.topic, "");
    }
    this.hint = "";
    this.topic = "";
  }
}
