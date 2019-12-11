import {Component} from "@angular/core";
import {MqttService} from "ngx-mqtt";
import {Message} from "../../message";
import {AppComponent} from "../../app.component";

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
    if (this.hint !== "" && this.hint !== undefined ) {
      this.app.sendInstruction([{
        instruction: "hint",
        value: this.hint
      }]);
      this.hint = "";
    }
  }
}
