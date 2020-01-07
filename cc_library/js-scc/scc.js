// var mqtt = require("mqtt");
// var client = mqtt.connect("mqtt://192.168.178.82");
//
// client.on("connect", onConnect);
// client.on("message", onMessage);
// function onConnect() {
//   client.subscribe("test", function(err) {
//     if (!err) {
//       console.log("started")
//       publish("test", "I'm alive")
//     } else {
//       console.log(err)
//     }
//   });
// }
//
// function onMessage(topic, message) {
//   console.log(topic + " : " + message.toString());
//   client.end();
// }
//
//
// function publish(topic, message) {
//     client.publish(topic, message)
// }
//
// publish("test", "hello sccLib");
//
// function test() {
//   console.log("TEST!!!")
// }
//
var SCC = (function SCC() {
  return function SCCConstruction() {
    var _this = this;
    _this.test = function () {
      console.log("TEST!!!")
    };
  };
}());

var scc = new SCC();
