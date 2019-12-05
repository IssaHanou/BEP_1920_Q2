import {Component, OnDestroy, OnInit, Type} from "@angular/core";
import {MqttService} from "ngx-mqtt";
import {Message} from "./message";
import {JsonConvert} from "json2typescript";
import {stringify} from "querystring";
import {MatSnackBar} from '@angular/material';
import {ConnectionStatus} from "ngx-mqtt-client";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.css", "../assets/css/main.css"]
})
export class AppComponent implements OnInit, OnDestroy {
  title = "S.C.I.L.E.R";
  nameOfRoom = "Super awesome escape";
  jsonConvert: JsonConvert;

  constructor(private _mqttService: MqttService, private snackBar: MatSnackBar) {
    this.jsonConvert = new JsonConvert()

    // /**
    //  * Tracks connection status.
    //  */
    // this._mqttService.status().subscribe((s: ConnectionStatus) => {
    //   const status = s === ConnectionStatus.CONNECTED ? 'CONNECTED' : 'DISCONNECTED';
    //   this.status.push(`Mqtt client connection status: ${status}`);
    // });
  }

  ngOnInit(): void {}

  /**
   * The purpose of this is, when the user leave the app we should cleanup our subscriptions
   * and close the connection with the broker
   */
  ngOnDestroy(): void {
    this._mqttService.disconnect();
  }

  // /**
  //  * Manages connection manually.
  //  * If there is an active connection this will forcefully disconnect that first.
  //  * @param {IClientOptions} config
  //  */
  // connect(config: IClientOptions): void {
  //   this._mqttService.connect(config);
  // }
  //
  // /**
  //  * Subscribes to fooBar topic.
  //  * The first emitted value will be a {@see SubscriptionGrant} to confirm your subscription was successful.
  //  * After that the subscription will only emit new value if someone publishes into the fooBar topic.
  //  * */
  // subscribe(): void {
  //   this._mqttService.subscribeTo<Foo>('fooBar')
  //     .subscribe({
  //       next: (msg: SubscriptionGrant | Foo) => {
  //         if (msg instanceof SubscriptionGrant) {
  //           this.status.push('Subscribed to fooBar topic!');
  //         } else {
  //           this.messages.push(msg);
  //         }
  //       },
  //       error: (error: Error) => {
  //         this.status.push(`Something went wrong: ${error.message}`);
  //       }
  //     });
  // }
  //
  //
  // /**
  //  * Sends message to fooBar topic.
  //  */
  // sendMsg(): void {
  //   this._mqttService.publishTo<Foo>('fooBar', {bar: 'foo'}).subscribe({
  //     next: () => {
  //       this.status.push('Message sent to fooBar topic');
  //     },
  //     error: (error: Error) => {
  //       this.status.push(`Something went wrong: ${error.message}`);
  //     }
  //   });
  // }
  //
  // /**
  //  * Unsubscribe from fooBar topic.
  //  */
  // unsubscribe(): void {
  //   this._mqttService.unsubscribeFrom('fooBar').subscribe({
  //     next: () => {
  //       this.status.push('Unsubscribe from fooBar topic');
  //     },
  //     error: (error: Error) => {
  //       this.status.push(`Something went wrong: ${error.message}`);
  //     }
  //   });
  // }
  //
  //




  /**
   * Send an instruction to the broker, over instruction topic.
   * @param instruction instruction to be sent.
   */
  public sendInstruction(instruction: string) {
    let msg = new Message("front-end",
      "instruction",
      new Date(),
      {"instruction": instruction}
    );
    let jsonMessage: string = this.jsonConvert.serialize(msg)
    this._mqttService.unsafePublish("instruction", JSON.stringify(jsonMessage));
    console.log("log: sent instruction message: " + JSON.stringify(jsonMessage))
  }

  /**
   * Process incoming message.
   * @param jsonMessage json message.
   */
  public processMessage(jsonMessage: string) {
    let msg: Message = Message.deserialize(jsonMessage);
    switch (msg.type) {
      case "confirmation": {
        // When the front-end receives confirmation message from client computer
        // that instruction was completed, show the message to the user.
        let display = "received confirmation from " + msg.deviceId
          + " for instruction " + stringify(msg.contents["instruction"]);
        this.snackBar.open(display);
        break;
      }
      case "instruction": {
        //TODO instructions to front-end? e.g. ask for hint
        break;
      }
      default:
        console.log("invalid instruction type " + msg.type);
        break;
    }
  }

}
