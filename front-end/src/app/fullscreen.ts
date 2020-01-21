/**
 * Class that can be extended to implement full screen abilities.
 * The setFullScreen method will be called by a button and determines whether to close or open accordingly.
 * This class uses the HTML5 full screen API.
 * Currently extended by app and camera components.
 */
export class FullScreen {
  fullScreen = false;

  constructor() {}

  /**
   * Open the full element as full screen.
   * @element the element to put into full screen.
   */
  openFullScreen(element) {
    const docElmWithBrowsersFullScreenFunctions = element as HTMLElement & {
      mozRequestFullScreen(): Promise<void>;
      webkitRequestFullscreen(): Promise<void>;
      msRequestFullscreen(): Promise<void>;
    };

    if (docElmWithBrowsersFullScreenFunctions.requestFullscreen) {
      docElmWithBrowsersFullScreenFunctions.requestFullscreen();
    } else if (docElmWithBrowsersFullScreenFunctions.mozRequestFullScreen) {
      /* Firefox */
      docElmWithBrowsersFullScreenFunctions.mozRequestFullScreen();
    } else if (docElmWithBrowsersFullScreenFunctions.webkitRequestFullscreen) {
      /* Chrome, Safari and Opera */
      docElmWithBrowsersFullScreenFunctions.webkitRequestFullscreen();
    } else if (docElmWithBrowsersFullScreenFunctions.msRequestFullscreen) {
      /* IE/Edge */
      docElmWithBrowsersFullScreenFunctions.msRequestFullscreen();
    }
    this.fullScreen = true;
  }

  /**
   * Close the full screen, the whole document should no longer be full screen.
   */
  closeFullScreen() {
    const docWithBrowsersExitFunctions = document as Document & {
      mozCancelFullScreen(): Promise<void>;
      webkitExitFullscreen(): Promise<void>;
      msExitFullscreen(): Promise<void>;
    };

    if (docWithBrowsersExitFunctions.exitFullscreen) {
      docWithBrowsersExitFunctions.exitFullscreen();
    } else if (docWithBrowsersExitFunctions.mozCancelFullScreen) {
      /* Firefox */
      docWithBrowsersExitFunctions.mozCancelFullScreen();
    } else if (docWithBrowsersExitFunctions.webkitExitFullscreen) {
      /* Chrome, Safari and Opera */
      docWithBrowsersExitFunctions.webkitExitFullscreen();
    } else if (docWithBrowsersExitFunctions.msExitFullscreen) {
      /* IE/Edge */
      docWithBrowsersExitFunctions.msExitFullscreen();
    }
    this.fullScreen = false;
  }

  /**
   * Determine whether to open of close full screen.
   * The default open is the full document element, but this can be overridden to open a certain element in full screen.
   */
  setFullScreen() {
    this.checkFullScreenEsc();

    if (this.fullScreen) {
      this.closeFullScreen();
    } else {
      this.openFullScreen(document.documentElement);
    }
  }

  /**
   * The full screen can be opened with a button (which sets the global fullScreen attribute of this class).
   * It can be closed again with the button, but also with the Esc key, which does not affect the attribute.
   * So, if the full screen boolean was not reset, do it manually.
   */
  checkFullScreenEsc() {
    const docWithBrowsersExitFunctions = document as Document & {
      mozCancelFullScreen(): Promise<void>;
      webkitExitFullscreen(): Promise<void>;
      msExitFullscreen(): Promise<void>;
    };
    const docFullScreen =
      docWithBrowsersExitFunctions.fullscreenElement != null;
    if (docFullScreen !== this.fullScreen) {
      this.fullScreen = !this.fullScreen;
    }
  }
}
