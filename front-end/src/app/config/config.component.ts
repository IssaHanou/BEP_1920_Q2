import { Component, OnInit } from "@angular/core";
import { AppComponent } from "../app.component";
import {logger} from "codelyzer/util/logger";

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

  /**
   * Define listeners for the file reader.
   * On reading file, it should check config on back-end.
   * On error, it should add to error list and log the error.
   */
  constructor(private app: AppComponent) {
    this.errors = [];
    this.reader = new FileReader();
    this.reader.addEventListener("load", (e) => {
      this.data = e.target["result"];
      this.app.sendInstruction([{ instruction: "check config", config: JSON.parse(this.data)}]);
    });
    this.reader.addEventListener("error", (e) => {
      this.errors.push(e.target["result"]);
      logger.error("log: error while reading file")
    })
  }

  ngOnInit() {
  }

  /**
   * When file is submitted, call file reader.
   */
  checkFile(files: FileList) {
    const file = files.item(0);
    this.uploaded = "Succesfully uploaded file: " + file.name + "!";
    this.reader.readAsText(file, "UTF-8");
  }

  getErrors(): string[] {
    const list = [];
    for (const err of this.errors) {
      list.push(err);
    }
    return list;
  }
}
