import { Component, OnInit } from "@angular/core";
import { AppComponent } from "../../app.component";
import {Device} from "../device/device";
import {Button} from "./button";

@Component({
  selector: "app-manage",
  templateUrl: "./manage.component.html",
  styleUrls: ["./manage.component.css", "../../../assets/css/main.css"]
})
export class ManageComponent implements OnInit {

  constructor(private app: AppComponent) {}

  ngOnInit() {}

  getButtons() {
    const buttons: Button[] = [];
    for (const btn of this.app.manageButtons.all.values()) {
      buttons.push(btn);
    }
    return buttons;
  }

  onClickTestButton() {
    this.app.sendInstruction([{ instruction: "test all" }]);
  }

  onClickResetButton() {
    this.app.sendInstruction([{ instruction: "reset all" }]);
    this.app.sendConnection(true);
  }

  /**
   * When clicking a button in the front-end manage section,
   * send updated data to the back-end.
   * All buttons have a boolean type, only update the pressed button.
   * @param btnID the button that is pressed
   */
  onClickCustomButton(btnID) {
    const device = this.app.deviceList.getDevice("front-end");
    if (device != null) {
      const status = device.status;
      const statusMsg = {};
      for (const component of this.app.manageButtons.all.values()) {
        let statusToSet = status.get(component.id);
        if (btnID === component.id) {
          statusToSet = !statusToSet;
        }
        statusMsg[component.id] = statusToSet;
      }
      this.app.sendStatus(statusMsg);
    }
  }
}
