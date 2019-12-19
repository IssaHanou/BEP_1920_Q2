import { Component, OnInit } from "@angular/core";
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
    this.app.sendInstruction([{ instruction: "test all" }]);
  }

  onClickResetButton() {
    this.app.sendInstruction([{ instruction: "reset all" }]);
    this.app.sendConnection(true);
  }

  onClickStartButton() {
    const status = this.app.deviceList.getDevice("front-end").status;
    this.app.sendStatus(status.get("start") + 1, status.get("stop"));
  }

  onClickStopButton() {
    const status = this.app.deviceList.getDevice("front-end").status;
    this.app.sendStatus(status.get("start"), status.get("stop") + 1);
  }
}
