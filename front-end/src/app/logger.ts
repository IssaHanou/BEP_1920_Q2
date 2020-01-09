import { formatMS } from "./components/timer/timer.component";
/**
 * This logger will log the messages from the application.
 * It is currently implemented to log to the console, but this can be changed to log to a file on the server-side.
 */
export class Logger {

  constructor() {}

  /**
   * Log a message with a certain level of severity.
   * @param level: info, warning or error
   * @param msg to log
   */
  public log(level: string, msg: string) {
    console.log("time=" + new Date().toLocaleString()
      + ", level=" + level
      + ", msg=" + msg);
  }
}
