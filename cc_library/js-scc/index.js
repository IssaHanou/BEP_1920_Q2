const Paho = require("paho-mqtt");

/**
 * Message is a class containing all info required by the message_manual.md
 */
class Message {
  constructor(device_id, type, contents) {
    this.device_id = device_id;
    this.time_sent = formatDate(new Date());
    this.type = type;
    this.contents = contents;
  }
}

/**
 * Abstract device class from which all custom devices should inherit.
 * Defines all required methods needer for communication to S.C.I.L.E.R.
 * @Abstract
 */
class Device {
  constructor(config, logger) {
    this.scclib = new SccLib(config, this, logger);

    // make sure abstract class Device cannot be instantiated directly
    if (this.constructor === Device) {
      throw new TypeError(
        "abstract class Device cannot be instantiated directly"
      );
    }

    // make sure abstract method getStatus is implemented when extending from Device
    if (typeof this.getStatus != "function") {
      throw new TypeError("abstract method 'getStatus' not implemented");
    }
    // make sure abstract method performInstruction is implemented when extending from Device
    if (typeof this.performInstruction != "function") {
      throw new TypeError(
        "abstract method 'performInstruction' not implemented"
      );
    }
    // make sure abstract method test is implemented when extending from Device
    if (typeof this.test != "function") {
      throw new TypeError("abstract method 'test' not implemented");
    }
    // make sure abstract method reset is implemented when extending from Device
    if (typeof this.reset != "function") {
      throw new TypeError("abstract method 'reset' not implemented");
    }
  }

  /**
   * statusChanged should be called whenever the status of a device changes
   */
  statusChanged() {
    this.scclib.statusChanged();
  }

  /**
   * log can be used to log in the same logger as this library
   * @param level one of the following strings: 'debug', 'info', 'warn', 'error', 'fatal'
   * @param message custom string containing more information
   */
  log(level, message) {
    this.scclib.log(level, message);
  }

  /**
   * start starts the device
   */
  start() {
    this.scclib.connect();
  }
}

/**
 * Class SccLib sets up the connection and the right handler
 */
class SccLib {
  constructor(config, device, logger) {
    // type check config
    const configProperties = ["id", "description", "host", "port", "labels"];
    for (let configProperty of configProperties) {
      if (!config.hasOwnProperty(configProperty)) {
        throw new TypeError(
          config + " should have a property: " + configProperty
        );
      }
    }

    // type check device
    if (device.prototype instanceof Device) {
      throw new TypeError(device + " should be of type Device");
    }

    // type check logger
    if (typeof logger !== "function") {
      throw new TypeError(
        logger + " should be of type function(date, level, message)"
      );
    }

    this.device = device;
    this.name = config.id;
    this.info = config.description;
    this.host = config.host;
    this.port = config.port;
    this.labels = config.labels;

    /**
     * log uses the log function provided to log level, time and a message
     * @param level one of the following strings: 'debug', 'info', 'warn', 'error', 'fatal'
     * @param message custom string containing more information
     */
    this.log = function(level, message) {
      logger(new Date(), level, message);
    };
    this.log("info", "Start of log for device: " + this.name);

    // setup mqtt
    this.client = new Paho.Client(this.host, this.port, "", this.name);
    this.client.onMessageArrived = msg => {
      this._onMessage(msg);
    };

    /**
     * _onConnect gets called when trying to connect,
     * it subscribes to all specified topics
     * it sends a connection true message
     * it logs that it connected
     * @private
     */
    this._onConnect = function() {
      // subscripe to all labels and standard topics
      for (let label in this.labels) {
        this.client.subscribe(label);
      }
      this.client.subscribe("client-computers");
      this.client.subscribe(this.name);

      // send connection status
      this._sendMessage(
        "back-end",
        new Message(this.name, "connection", {
          connection: true
        })
      );

      // log successful connection
      this.log("info", "connected OK");
    };

    /**
     * _onConnectFailure is a method that gets called when connecting failed
     * it will log an error and try to reconnect on a regular interval till it succeeds
     * @private
     */
    this._onConnectFailure = function() {
      let retryCooldown = 10 * 1000; // 10 seconds before retrying to connect
      this.log(
        "error",
        "connecting failed, retry in " + retryCooldown + " seconds"
      );
      setTimeout(() => {
        this.connect();
      }, retryCooldown);
    };

    this._onMessage = function(message) {
      this.log(
        "info",
        "message received:\n topic: " +
          message.topic +
          ",\n message: " +
          message.payloadString
      );
      this._handle(message.payloadString);
    };
  }

  _handle(payloadString) {
    const message = JSON.parse(payloadString);
    if (message.type !== "instruction") {
      this.log(
        "warn",
        "received non-instruction message of type: " + message.type
      );
    } else {
      const success = this._checkMessage(message.contents);
    }
  }

  _checkMessage(contents) {
    for (let i = 0; i < contents.length; i++) {
      const action = contents[i];

      const instruction = action.instruction;
      switch (instruction) {
        case "test": {
          this.device.test();
          this.log("info", "instruction performed " + instruction);
          break;
        }
        case "status update": {
          const message = new Message(this.name, "connection", {
            connection: true
          });
          this._sendMessage("back-end", message);
          this.statusChanged();
          this.log("info", "instruction performed " + instruction);
          break;
        }
        case "reset": {
          this.device.reset();
          this.log("info", "instruction performed " + instruction);
          break;
        }
        default: {
          break;
        }
      }
    }
  }

  /**
   * _sendMessage sends an mqtt message to SCILER
   * @param topic string containing the mqtt topic
   * @param message json that should follow the message_manual.md
   * @private
   */
  _sendMessage(topic, message) {
    let msg = new Paho.Message(JSON.stringify(message));
    msg.destinationName = topic;
    this.client.send(msg);
  }

  /**
   * connect connects to SCILER
   * sets up a LWT (Last Will and Testament)
   * sets up handlers for connection and connection failure
   * sets up automatic reconnect
   */
  connect() {
    let will = new Paho.Message(
      JSON.stringify(
        new Message(this.name, "connection", {
          connection: false
        })
      )
    );
    will.destinationName = "back-end";
    this.client.connect({
      onSuccess: () => {
        this._onConnect();
      },
      onFailure: () => {
        this._onConnectFailure();
      },
      willMessage: will,
      reconnect: true,
      keepAliveInterval: 10
    });
  }

  statusChanged() {
    this._sendMessage(
      "back-end",
      new Message(this.name, "status", this.device.getStatus())
    );
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

module.exports = Device;
