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

    this.client = new Paho.Client(this.host, this.port, "", this.name);
    this.client.onMessageArrived = msg => {
      this._onMessage(msg);
    };

    // this.connect = function () {
    //   this.client.connect({
    //     onSuccess: () => {
    //       this._onConnect();
    //     },
    //     onFailure: err => {
    //       this._onConnectFailure(err);
    //     } // todo wait and try to connect again
    //   });
    // };

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

    this._onMessage = function(message) {
      this.log(
        "info",
        "message received:\n topic: " +
          message.topic +
          ",\n message: " +
          message.payloadString
      );
    };
  }

  connect() {
    let will = new Paho.Message(
      JSON.stringify({
        topic: "back-end",
        payloadString: JSON.stringify({
          device_id: this.name,
          type: "connection",
          timeSent: formatDate(new Date()),
          contents: {
            connection: false
          }
        })
      })
    );
    will.destinationName = "back-end";
    this.client.connect({
      onSuccess: () => {
        this._onConnect();
      },
      onFailure: err => {
        this._onConnectFailure(err);
      }, // todo wait and try to connect again
      willMessage: will
    });
  }
}

const formatDate = function(date) {
  return (
    date.getDate() +
    "-" +
    date.getMonth() +
    1 +
    "-" +
    date.getFullYear() +
    " " +
    date.getHours() +
    ":" +
    date.getMinutes() +
    ":" +
    date.getSeconds()
  );
};
module.exports = SccLib;
