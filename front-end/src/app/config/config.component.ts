import { Component, OnInit } from "@angular/core";
import { AppComponent } from "../app.component";

/**
 * The config component controls the configuration page, selected through the side menu.
 */
@Component({
  selector: "app-config",
  templateUrl: "./config.component.html",
  styleUrls: ["./config.component.css"]
})
export class ConfigComponent implements OnInit {
  reader: FileReader;
  data = "";
  errors: string[];
  currentFile: File;

  /**
   * Define listeners for the file reader, which is used in the file reading button.
   * First reset the errors.
   * On reading file, it should check config on back-end.
   * On error, it should add to error list and log the error.
   */
  constructor(private app: AppComponent) {
    this.errors = [];
    this.reader = new FileReader();
    const res = "result";
    this.reader.addEventListener("load", e => {
      this.data = e.target[res];
      this.errors = [];
      this.app.sendInstruction([
        {
          instruction: "check config",
          config: this.getJSONConfig(),
          name: this.currentFile.name
        }
      ]);
    });
    this.reader.addEventListener("error", e => {
      this.errors.push("level I - file error: " + e.target[res]);
      this.app.logger.log("error", "error while reading file");
      this.app.uploadedConfig = "Error tijdens uploaden: " + e.target[res];
    });
  }

  ngOnInit() {}

  /**
   * When clicking the upload button, first reset the file, so the same one can be entered again for re-checking.
   */
  resetCurrentFile(event) {
    event.target.value = null;
  }

  /**
   * When file is submitted, call the file reader and tell the user the upload was successful.
   * Only one file can be submitted.
   */
  checkFile(files: FileList) {
    this.currentFile = files.item(0);
    this.reader.readAsText(this.currentFile, "UTF-8");
  }

  /**
   * Use the config entered as new configuration for app, by messaging the back-end.
   */
  sendConfig() {
    this.app.resetConfig();
    this.app.sendInstruction([
      {
        instruction: "use config",
        config: this.getJSONConfig(),
        file: this.currentFile.name
      }
    ]);
  }

  /**
   * Parses JSON and catches error when JSON is invalid.
   */
  getJSONConfig() {
    try{
      return JSON.parse(this.data)
    } catch (e) {
      this.errors.push("level I - JSON error: "+ e);
      this.app.logger.log("error", "error while reading file");
      this.app.uploadedConfig = "Error tijdens uploaden: " + e;
    }
  }

  /**
   * All JSON errors will be shown in the box, but one on each line.
   */
  getErrors(): string[] {
    const list = [];
    for (const err of this.errors) {
      list.push(err);
    }
    for (const err of this.app.configErrorList) {
      list.push(err);
    }
    return list;
  }

  /**
   * Check whether no errors were found, both through reading and in the back-end.
   */
  noErrors(): boolean {
    return this.errors.length === 0 && this.app.configErrorList.length === 0;
  }

  /**
   * Return the name of the uploaded config file.
   */
  getUploaded(): string {
    return this.app.uploadedConfig;
  }
}
