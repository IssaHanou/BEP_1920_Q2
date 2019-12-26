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
  camera: string;

  constructor(private app: AppComponent, private sanitizer: DomSanitizer) {
  }

  ngOnInit() {
    this.cameraFeedSrc = this.sanitizer.bypassSecurityTrustResourceUrl("about:blank");
  }

  allCameras(): Camera[] {
    return this.app.cameras;
  }

  /**
   * Set the src of the iframe with parameter from dropdown.
   * Dom sanitizer used to make sure link from config is 'safe'.
   */
  setSrc() {
    for (const cam of this.allCameras()) {
      if (this.camera === cam.name) {
        this.cameraFeedSrc = this.sanitizer.bypassSecurityTrustResourceUrl(cam.link);
      }
    }
  }
}
