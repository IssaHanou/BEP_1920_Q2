import { Component, OnInit } from "@angular/core";
import { AppComponent } from "../../app.component";
import { Button } from "./button";

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
   * Return the buttons for the manage section, in alphabetical order.
   */
  getButtons() {
    const buttons: Button[] = [];
    for (const btn of this.app.manageButtons.all.values()) {
      buttons.push(btn);
    }
    buttons.sort((a, b) => a.id.localeCompare(b.id));
    return buttons;
  }

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
