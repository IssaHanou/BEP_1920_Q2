import {Component, OnDestroy, OnInit} from "@angular/core";
import { AppComponent } from "../app.component";
import { Camera } from "./camera";
import { DomSanitizer } from "@angular/platform-browser";
import { FormControl } from "@angular/forms";

@Component({
  selector: "app-camera",
  templateUrl: "./camera.component.html",
  styleUrls: ["./camera.component.css"]
})
export class CameraComponent implements OnInit, OnDestroy {

  selectedCameraControl = new FormControl();
  cameraFeedSrc: any;
  selectedCameraControl2 = new FormControl();
  cameraFeedSrc2: any;
  openSecond = false;

  /**
   * The app is used for the list of all cameras.
   * It also keeps track of the selected cameras and open second variable,
   * to keep them on when returning to the camera page after switching.
   */
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
    if (this.app.selectedCamera2 !== undefined) {
      this.selectedCameraControl2.setValue(this.app.selectedCamera2);
      this.setSrc2();
    } else if (this.allCameras().length > 1) {
      this.selectedCameraControl2.setValue(this.allCameras()[0].name);
      this.setSrc2();
    } else {
      this.cameraFeedSrc2 = this.sanitizer.bypassSecurityTrustResourceUrl("about:blank");
    }

    this.openSecond = this.app.openSecondCamera;
  }

  ngOnDestroy(): void {
    this.app.openSecondCamera = this.openSecond;
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

  /**
   * Set the src of the secondiframe with parameter from dropdown.
   * Dom sanitizer used to make sure link from config is 'safe'.
   */
  setSrc2() {
    for (const cam of this.allCameras()) {
      if (this.selectedCameraControl2.value === cam.name) {
        this.cameraFeedSrc2 = this.sanitizer.bypassSecurityTrustResourceUrl(cam.link);
        this.app.selectedCamera2 = this.selectedCameraControl2.value;
      }
    }
  }
}
