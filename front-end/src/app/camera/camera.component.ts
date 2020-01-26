import {
  AfterViewInit,
  Component,
  ElementRef,
  OnDestroy,
  OnInit,
  ViewChild
} from "@angular/core";
import { AppComponent } from "../app.component";
import { Camera } from "./camera";
import { DomSanitizer } from "@angular/platform-browser";
import { FormControl } from "@angular/forms";
import { FullScreen } from "../fullscreen";

/**
 * The camera components controls the camera page, selected through the side menu.
 * It allows for the box to become full screen.
 */
@Component({
  selector: "app-camera",
  templateUrl: "./camera.component.html",
  styleUrls: ["./camera.component.css"]
})
export class CameraComponent extends FullScreen
  implements OnInit, AfterViewInit, OnDestroy {
  /**
   * There is a control for each selection dropdown and the src used by each iframe.
   * The booleans keep track of the second feed (opened or not) and full screen (enabled or not).
   */
  selectedCameraControl = new FormControl();
  cameraFeedSrc: any;
  selectedCameraControl2 = new FormControl();
  cameraFeedSrc2: any;
  openSecond = false;
  fullScreen = false;

  /**
   * Keeping track of div boxes, for styling.
   */
  @ViewChild("cameraBox", { static: true }) cameraBox: ElementRef;
  @ViewChild("contents", { static: true }) boxContents: ElementRef;

  /**
   * The app is used for the list of all cameras.
   * It also keeps track of the selected cameras and open second variable,
   * to keep them on when returning to the camera page after switching.
   */
  constructor(private app: AppComponent, private sanitizer: DomSanitizer) {
    super();
  }

  /**
   * Upon initializing, make sure to set the feeds of the iframe's to either:
   * 1. previously selected value
   * 2. the first value in the list
   * 3. a blank status to prevent crashes
   */
  ngOnInit() {
    if (this.app.selectedCamera !== undefined) {
      this.selectedCameraControl.setValue(this.app.selectedCamera);
      this.setSrc();
    } else if (this.allCameras().length > 1) {
      this.selectedCameraControl.setValue(this.allCameras()[0].name);
      this.setSrc();
    } else {
      this.cameraFeedSrc = this.sanitizer.bypassSecurityTrustResourceUrl(
        "about:blank"
      );
    }
    if (this.app.selectedCamera2 !== undefined) {
      this.selectedCameraControl2.setValue(this.app.selectedCamera2);
      this.setSrc2();
    } else if (this.allCameras().length > 1) {
      this.selectedCameraControl2.setValue(this.allCameras()[0].name);
      this.setSrc2();
    } else {
      this.cameraFeedSrc2 = this.sanitizer.bypassSecurityTrustResourceUrl(
        "about:blank"
      );
    }
    this.openSecond = this.app.openSecondCamera;
  }

  ngAfterViewInit(): void {}

  /**
   * When switching screens in the menu, have the application keep track of whether or not a second screen was opened.
   */
  ngOnDestroy(): void {
    this.app.openSecondCamera = this.openSecond;
  }

  /**
   * Returns the list of cameras.
   */
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
        this.cameraFeedSrc = this.sanitizer.bypassSecurityTrustResourceUrl(
          cam.link
        );
        this.app.selectedCamera = this.selectedCameraControl.value;
      }
    }
  }

  /**
   * Set the src of the second iframe with parameter from dropdown.
   * Dom sanitizer used to make sure link from config is 'safe'.
   */
  setSrc2() {
    for (const cam of this.allCameras()) {
      if (this.selectedCameraControl2.value === cam.name) {
        this.cameraFeedSrc2 = this.sanitizer.bypassSecurityTrustResourceUrl(
          cam.link
        );
        this.app.selectedCamera2 = this.selectedCameraControl2.value;
      }
    }
  }

  /**
   * When the camera is opened in full screen, disable the header expansion panel for selecting other feeds.
   */
  expandPanel(matExpansionPanel, event): void {
    event.stopPropagation();

    if (this.fullScreen) {
      matExpansionPanel.close();
    }
  }

  /**
   * Sets height of contents box to full and calls super function, to the full screen of camera box.
   */
  openFullScreen() {
    this.boxContents.nativeElement.style.height = "90vh";
    super.openFullScreen(this.cameraBox.nativeElement);
    this.fullScreen = true;
  }

  /**
   * Sets height of contents box back to original and calls super function.
   */
  closeFullScreen() {
    // Trigger fullscreen
    this.boxContents.nativeElement.style.height = "70vh";
    super.closeFullScreen();
    this.fullScreen = false;
  }

  /**
   * Calls the appropriate close and open method for full screen.
   * Overrides super method.
   */
  setFullScreen() {
    this.checkFullScreenEsc();

    if (this.fullScreen) {
      this.closeFullScreen();
    } else {
      this.openFullScreen();
    }
  }
}
