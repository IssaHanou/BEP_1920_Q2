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
class SccLib {
  constructor(config, device) {
    this.device = device;
    // config = JSON.parse(config);
    this.name = config.id;
    this.info = config.description;
    this.host = config.host;
    this.port = config.port;
    this.labels = config.labels
  }

  test() {
    console.log("name: " + this.name + " " + this.info);
  }
}
export { SccLib };
