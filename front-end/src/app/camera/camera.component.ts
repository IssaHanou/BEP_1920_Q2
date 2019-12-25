import { Component, OnInit } from "@angular/core";
import { AppComponent } from "../app.component";
import { Camera } from "./camera";
import { DomSanitizer } from "@angular/platform-browser";

@Component({
  selector: "app-camera",
  templateUrl: "./camera.component.html",
  styleUrls: ["./camera.component.css"]
})
export class CameraComponent implements OnInit {

  cameraFeedSrc: any;

  constructor(private app: AppComponent, private sanitizer: DomSanitizer) {
  }

  ngOnInit() {
    this.setSrc();
  }

  setSrc() {
    this.cameraFeedSrc = this.sanitizer.bypassSecurityTrustResourceUrl("https://raccoon.games");
  }

  allCameras(): Camera[] {
    return this.app.cameras;
  }
}
