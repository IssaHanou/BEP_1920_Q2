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
    buttons.sort((a, b) => a.id.localeCompare(b.id));
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
   * update the status of clicked button and
   * send updated data to the back-end.
   * @param btnID the button that is pressed
   */
  onClickCustomButton(btnID) {
    this.app.deviceList.updateDevice(btnID, true);
    this.app.sendStatusFrontEnd();

    this.app.deviceList.updateDevice(btnID, false);
    this.app.sendStatusFrontEnd();
  }
}
