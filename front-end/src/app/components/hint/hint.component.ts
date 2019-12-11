import { Component } from "@angular/core";
import { AppComponent } from "../../app.component";

@Component({
  selector: "app-hint",
  templateUrl: "./hint.component.html",
  styleUrls: ["./hint.component.css", "../../../assets/css/main.css"]
})
export class HintComponent {
  hint: string;

  constructor(private app: AppComponent) {
  }

  onSubmit() {
    this.app.sendInstruction([{
      instruction: "hint",
      value: this.hint
    }]);
    this.hint = "";
  }
}
