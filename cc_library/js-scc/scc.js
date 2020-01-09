const Paho = require("paho-mqtt");

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

    this.client = new Paho.Client("localhost", 8083, "", "fancy-display");
    this.client.onMessageArrived = (msg) => {this._onMessage(msg)};

    this.client.connect({
      onSuccess: () => {this._onConnect()},
      onFailure: (err) => {this._onConnectFailure(err)}, // todo wait and try to connect again
    });

    this._onConnect = function() {
      this.client.subscribe("test");
      let message = new Paho.Message("KAAAAAS!!!");
      message.destinationName = "test";
      this.client.send(message);
      this.log("info", "Connected!");
    };

    this._onConnectFailure = function(err) {
      this.log("error", err.errorMessage);
    };

    this._onMessage = function (message) {
      this.log("info", "msg: " + message.payloadString)
    };
  }
}
export { SccLib };
