import { Component, OnInit } from "@angular/core";
import { AppComponent } from "../app.component";
import { Camera } from "./camera";
import { DomSanitizer } from "@angular/platform-browser";
import { FormControl } from "@angular/forms";

@Component({
  selector: "app-camera",
  templateUrl: "./camera.component.html",
  styleUrls: ["./camera.component.css"]
})
export class CameraComponent implements OnInit {

  selectedCameraControl = new FormControl();
  cameraFeedSrc: any;

  constructor(private app: AppComponent, private sanitizer: DomSanitizer) {
  }

  ngOnInit() {
    if (this.app.selectedCamera !== undefined) {
      this.selectedCameraControl.setValue(this.app.selectedCamera);
      this.setSrc();
    } else if (this.allCameras().length > 1) {
      this.selectedCameraControl.setValue(this.allCameras()[0].name);
      this.setSrc();
    } else {
      this.cameraFeedSrc = this.sanitizer.bypassSecurityTrustResourceUrl("about:blank");
    }
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
      if (this.selectedCameraControl.value === cam.name) {
        this.cameraFeedSrc = this.sanitizer.bypassSecurityTrustResourceUrl(cam.link);
        this.app.selectedCamera = this.selectedCameraControl.value;
      }
    }
  }
}
