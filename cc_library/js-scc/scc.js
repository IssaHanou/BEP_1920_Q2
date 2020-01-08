const mqtt = require("mqtt")

class SccLib {
  constructor(config, device, logger) {
    this.device = device;
    this.name = config.id;
    this.info = config.description;
    this.host = config.host;
    this.port = config.port;
    this.labels = config.labels;
    this.log = function(level, message) {
      logger(new Date(), level, message);
    };
    this.log("info", "Start of log for device: " + this.name);

    this._onConnect = function() {
      var client = this;
      client.subscribe("test", function(err) {
        if (!err) {
          console.log("started");
          client.publish("test", "I'm alive");
        } else {
          console.log(err);
        }
      });
    };

    this._onMessage = function(topic, message) {
      console.log(topic + " : " + message.toString());
      this.end();
    };

    this.publish = function (topic, message) {
      this.client.publish(topic, message);
      console.log("published! cl");
      this.log("debug", "published ll");
    };

    // MQTT

    this.client = mqtt.connect("ws://" + this.host, {
      port: this.port,
      clientId: this.name
    });
    this.client.on("connect", this._onConnect);
    this.client.on("message", this._onMessage);
  }
}
export { SccLib };
// var sccLib = new SccLib(
//   {
//     id: "device",
//     description: "bla bla bla",
//     host: "localhost",
//     port: 8083,
//     labels: []
//   },
//   "",
//   function(time, level, message) {
//     console.log(message);
//   }
// );
