import { Component } from "@angular/core";
import { AppComponent } from "../../app.component";
import { Hint } from "./hint";

@Component({
  selector: "app-hint",
  templateUrl: "./hint.component.html",
  styleUrls: ["./hint.component.css", "../../../assets/css/main.css"]
})
export class HintComponent {
  customHint: string;
  predefinedHint: string;

  constructor(private app: AppComponent) {}

  getPuzzleList(): Hint[] {
    const list = [];
    for (const hint of this.app.hintList) {
      list.push(hint);
    }
    return list;
  }

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

  onPredefinedHint() {
    if (this.predefinedHint !== undefined && this.predefinedHint !== "" && this.predefinedHint !== "---") {
      this.app.sendInstruction([
        {instruction: "hint", value: this.predefinedHint}
      ]);
      this.predefinedHint = "---";
    }
  }

  onCustomHint() {
    if (this.customHint !== undefined && this.customHint !== "") {
      this.app.sendInstruction([
        {
          instruction: "hint",
          value: this.customHint
        }
      ]);
      this.customHint = "";
    }
  }
}
