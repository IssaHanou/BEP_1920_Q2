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
  log(level: string, msg: string) {
    console.log ("time=" + new Date()
      + ", level=" + level
      + ", msg=" + msg);
  }
}
