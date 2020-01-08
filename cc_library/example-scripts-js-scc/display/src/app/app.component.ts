import { Component, OnInit } from "@angular/core";
import { SccLib } from "../../../../js-scc/scc"; // can also be replace with import to node and npm link
import { HttpClient } from "@angular/common/http";
import * as Paho from "paho-mqtt";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.css"]
})
export class AppComponent implements OnInit {
  title = "display";
  hint = "";
  scc;

  constructor(private http: HttpClient) {}

  ngOnInit(): void {
    // this.http
    //   .get("assets/display_config.json")
    //   .toPromise()
    //   .then((response: any) => {
    //     const config = response;
    //     this.scc = new SccLib(config, 4, function(date, level, message) {
    //       const formatDate = function(date) {
    //         return (
    //           date.getDate() +
    //           "-" +
    //           date.getMonth() +
    //           1 +
    //           "-" +
    //           date.getFullYear() +
    //           " " +
    //           date.getHours() +
    //           ":" +
    //           date.getMinutes() +
    //           ":" +
    //           date.getSeconds()
    //         );
    //       };
    //       console.log(
    //         "time=" + formatDate(date) + " level=" + level + " msg=" + message
    //       ); // call own logger
    //     });
    //     // this.scc.publish("test", "hoi hoi");
    //   });

    var client = new Paho.Client("localhost", 8083, "fancy-display");
    client.connect({onSuccess:_onConnect});


    function _onConnect() {
      client.subscribe("test");
      var message = new Paho.Message("KAAAAAS!!!");
      message.destinationName = "test";
      client.send(message);
      console.log("Connected!");
    }
  }
}
