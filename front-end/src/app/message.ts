import { JsonObject, JsonProperty } from "json2typescript";
import * as moment from "moment";

/**
 * Message class with same variables as message json object that is sent by broker.
 */
@JsonObject
export class Message {
  @JsonProperty("device_id")
  deviceId: string;
  @JsonProperty("type")
  type: string;
  @JsonProperty("time_sent")
  timeSent: string;
  @JsonProperty("contents")
  contents: any;

  constructor(deviceId: string, type: string, date: Date, contents: any) {
    this.deviceId = deviceId;
    this.timeSent = Message.formatDate(date);
    this.type = type;
    this.contents = contents;
  }

  /**
   * Format date to proper structure: dd-mm-yyyy hh:mm:ss
   * @param date Date object
   */
  private static formatDate(date: Date): string {
    date.setMonth(date.getMonth() - 1);
    return moment(date).format("DD-MM-YYYY HH:mm:ss");
  }

  /**
   * Transform json message into Message object.
   * @param jsonMessage incoming json message
   */
  public static deserialize(jsonMessage: string): Message {
    const msg = JSON.parse(jsonMessage);
    const deviceId = msg.device_id;
    const timeSent = msg.time_sent;
    const dateTime = timeSent.split(" ");
    const date = dateTime[0].split("-");
    const time = dateTime[1].split(":");
    const newDate = new Date(
      date[2],
      date[1],
      date[0],
      time[0],
      time[1],
      time[2]
    );
    const type = msg.type;
    const contents = msg.contents;
    return new Message(deviceId, type, newDate, contents);
  }
}
