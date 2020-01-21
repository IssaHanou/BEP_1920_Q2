export class FullScreen {
  fullScreen = false;

  constructor() {}

  /**
   * Uses HTML5 API to open the full screen.
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
   * Uses HTML5 API to close the full screen.
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
   * When the full screen is opened with the button, but closed with Esc,
   * the full screen boolean is not reset, so do it manually.
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
