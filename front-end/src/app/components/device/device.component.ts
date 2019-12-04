import {Component, OnDestroy, OnInit} from "@angular/core";
import {Subscription} from "rxjs";
import {IMqttMessage, MqttService} from "ngx-mqtt";

@Component({
  selector: "app-device",
  templateUrl: "./device.component.html",
  styleUrls: ["./device.component.css", "../../../assets/css/main.css"]
})


export class DeviceComponent implements OnInit {
  jsonData: JSON
  jsonMsg: JSON
  msg: IMqttMessage
  private subscription: Subscription;
  topicname: "status";

  onMessageArrived(message) {
    if (message.destinationName.indexOf("status") !== -1) {
      this.jsonMsg = JSON.parse(message.payloadString);
    }


  }

  subscribeNewTopic(): void {
    console.log("inside subscribe new topic")
    this.subscription = this.mqttService.observe(this.topicname).subscribe((message: IMqttMessage) => {
      this.msg = message;
      console.log("Message: " + message.payload.toString() + "<br> for topic: " + message.topic)
      this.onMessageArrived(message)
    });
    console.log("subscribed to topic: " + this.topicname)
  }


  constructor(private mqttService: MqttService) {
  }

  ngOnInit() {
    const column = ["apparaat", "connectie", "status"]
    const data = `[
        {
          "id": "Door",
          "status": [
                       {
                        "id": "door",
                        "status": true
                       }
                    ],
           "connection":"true"
        },
        {
          "id": "Buttons",
          "status": [
                    {
                      "id": "knop3",
                       "status": true
                    },
                    {
                      "id": "knop5",
                      "status": false
                    }
                    ],
          "connection": "true"
        }
    ]`;

    this.jsonData = JSON.parse(data)
  }
}
