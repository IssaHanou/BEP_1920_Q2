import { Component, OnInit } from "@angular/core";
import { AppComponent } from "../app.component";

@Component({
  selector: "app-config",
  templateUrl: "./config.component.html",
  styleUrls: ["./config.component.css"]
})
export class ConfigComponent implements OnInit {

  uploaded: string = "";
  data: string = "";
  errors: string[];
  reader: FileReader;
  noErrors: boolean;
  newConfig: string;
  currentFile: File;

  /**
   * Define listeners for the file reader.
   * First resets the errors.
   * On reading file, it should check config on back-end.
   * On error, it should add to error list and log the error.
   */
  constructor(private app: AppComponent) {
    this.errors = [];
    this.reader = new FileReader();
    this.reader.addEventListener("load", (e) => {
      this.data = e.target["result"];
      this.errors = [];
      this.app.configErrorList = [];
      this.app.sendInstruction([{ instruction: "check config", config: this.data}]);
    });
    this.reader.addEventListener("error", (e) => {
      this.errors.push(e.target["result"]);
      console.log("log: error while reading file");
      this.uploaded = "Error tijdens uploaden: " + e.target["result"];
    });
  }

  ngOnInit() {
  }

  /**
   * When file is submitted, call file reader.
   */
  checkFile(files: FileList) {
    this.currentFile = files.item(0);
    this.uploaded = "Uploaden gelukt: " + this.currentFile.name + "!";
    this.reader.readAsText(this.currentFile, "UTF-8");
    this.noErrors = this.getErrors().length === 0;
  }

  /**
   * Use the config entered as new configuration for app.
   */
  sendConfig() {
    this.app.sendInstruction([{instruction: "use config", config: this.data}])
    this.newConfig = "Configuratie uit: " + this.currentFile.name + " wordt nu gebruikt";
  }

  getErrors(): string[] {
    const list = [];
    for (const err of this.app.configErrorList) {
      list.push(err);
    }
    for (const err of this.errors) {
      list.push(err)
    }
    return list;
  }
}
