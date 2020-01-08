const Paho = require("paho-mqtt")

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

    this.client = new Paho.Client(this.host, 8083, this.name);
    console.log(this.client);
    this.client.connect({onSuccess:_onConnect, onFailure: _onConnectFailure});


    function _onConnect() {
      this.client.subscribe("test");
      var message = new Paho.Message("KAAAAAS!!!");
      message.destinationName = "test";
      this.client.send(message);
      this.log("info", "Connected!");
    }

    function _onConnectFailure(err) {
      this.log("error", err.errorMessage)
    }

    this.log("debug", "TEST!")

    // this._onMessage = function(topic, message) {
    //   console.log(topic + " : " + message.toString());
    //   this.end();
    // };
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
