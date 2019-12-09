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
    this.app.sendInstruction("test all");
  }

  /**
   * Test for the processing of messages, for now just a placeholder for confirming start instruction.
   */
  onClickStartButton() {
    const msg = new Message("front-end", "confirmation", new Date(), {
      completed: true,
      instructed: {
        device_id: "door",
        time_sent: "10-05-2019 15:09:14",
        type: "instruction",
        contents: { instruction: "start" }
      }
    });
    const res = this.app.jsonConvert.serialize(msg);
    this.app.processMessage(JSON.stringify(res));
  }
}
