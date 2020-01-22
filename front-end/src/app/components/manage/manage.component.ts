import { Component, OnInit } from "@angular/core";
import { AppComponent } from "../../app.component";

/**
 * The manage component controls the button section in the "Acties" box on the home page.
 */
@Component({
  selector: "app-manage",
  templateUrl: "./manage.component.html",
  styleUrls: ["./manage.component.css", "../../../assets/css/main.css"]
})
export class ManageComponent implements OnInit {
  constructor(private app: AppComponent) {}

  ngOnInit() {}

  /**
   * When test button is clicked, tell the back-end to let all devices test.
   */
  onClickTestButton() {
    this.app.sendInstruction([{ instruction: "test all" }]);
  }

  /**
   * When reset button is clicked, tell the back-end to let all devices reset.
   */
  onClickResetButton() {
    this.app.sendInstruction([{ instruction: "reset all" }]);
    this.app.sendConnection(true);
  }

  /**
   * When start button is clicked, tell the back-end to start the game.
   */
  onClickStartButton() {
    const device = this.app.deviceList.getDevice("front-end");
    if (device != null) {
      const status = device.status;
      this.app.sendStatus(status.get("start") + 1, status.get("stop"));
    }
  }

  /**
   * When start button is clicked, tell the back-end to stop the game.
   */
  onClickStopButton() {
    const device = this.app.deviceList.getDevice("front-end");
    if (device != null) {
      const status = device.status;
      this.app.sendStatus(status.get("start"), status.get("stop") + 1);
    }
  }
}
