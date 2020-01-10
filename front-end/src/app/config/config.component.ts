import { Component, OnInit } from "@angular/core";
import { AppComponent } from "../app.component";

@Component({
  selector: "app-config",
  templateUrl: "./config.component.html",
  styleUrls: ["./config.component.css"]
})
export class ConfigComponent implements OnInit {

  reader: FileReader;
  uploaded = "";
  data = "";
  errors: string[];
  currentFile: File;

  /**
   * Define listeners for the file reader, which is used in the file reading button.
   * First resets the errors.
   * On reading file, it should check config on back-end.
   * On error, it should add to error list and log the error.
   */
  constructor(private app: AppComponent) {
    this.errors = [];
    this.reader = new FileReader();
    const res = "result";
    this.reader.addEventListener("load", (e) => {
      this.data = e.target[res];
      this.errors = [];
      this.app.configErrorList = [];
      this.app.sendInstruction(
        [{
          instruction: "check config",
          config: JSON.parse(this.data)
        }]);
    });
    this.reader.addEventListener("error", (e) => {
      this.errors.push(e.target[res]);
      console.log("log: error while reading file");
      this.uploaded = "Error tijdens uploaden: " + e.target[res];
    });
  }

  ngOnInit() {
  }

  /**
   * When file is submitted, call file reader.
   */
  checkFile(files: FileList) {
    this.currentFile = files.item(0);
    this.reader.readAsText(this.currentFile, "UTF-8");
    this.uploaded = "Uploaden gelukt: " + this.currentFile.name + "!";
  }

  /**
   * Use the config entered as new configuration for app.
   */
  sendConfig() {
    this.app.sendInstruction([{
      instruction: "use config",
      config: JSON.parse(this.data),
      file: this.currentFile.name
    }]);
  }

  /**
   * All JSON errors will be shown per one - json unmarshal
   */
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

  noErrors(): boolean {
    return this.errors.length === 0 && this.app.configErrorList.length === 0;
  }
}
