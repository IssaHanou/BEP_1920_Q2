import { Component, OnInit } from "@angular/core";
import { Message } from "../../message";
import { AppComponent } from "../../app.component";

@Component({
  selector: "app-manage",
  templateUrl: "./manage.component.html",
  styleUrls: ["./manage.component.css", "../../../assets/css/main.css"]
})
export class ManageComponent implements OnInit {
  constructor(private app: AppComponent) {}

  ngOnInit() {}

  onClickTestButton() {
    this.app.sendInstruction("test all")
  }

  /**
   * Test for the processing of messages.
   * TODO delete this
   */
  onClickStartButton() {
    let msg = new Message("front-end",
      "instruction",
      new Date(),
      {"instruction": "test"}
    );
    let res = this.app.jsonConvert.serialize(msg);
    this.app.processMessage(JSON.stringify(res));
  }
}
