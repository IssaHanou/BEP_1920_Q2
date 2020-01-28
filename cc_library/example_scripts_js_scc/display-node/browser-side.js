$(document).ready(function() {
  const Device = require("sciler");
  let display;

  class Display extends Device {
    constructor(config) {
      super(config, timedLogger);
      this.hint = "";
      this.button = false;
    }

    // required method for extending Device
    getStatus() {
      return {
        button: this.button,
        hint: this.hint
      };
    }

    // required method for extending Device
    performInstruction(action) {
      switch (action.instruction) {
        case "hint": {
          this.hint = action.value;
          displayText(this.hint);
          this.statusChanged();
          break;
        }
        case "time": {
          displayText(action.duration)
        }
        default: {
          return false;
        }
      }
      return true;
    }

    // required method for extending Device
    test() {
      this.hint = "test";
      displayText(this.hint);
    }

    // required method for extending Device
    reset() {
      this.hint = "";
      this.button = false;
      displayText(this.hint);
      this.statusChanged();
    }
  }

  // custom logger used in constructor
  function timedLogger(date, level, message) {
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
    console.log(
      "time=" + formatDate(date) + " level=" + level + " msg=" + message
    ); // call own logger);
  }

  // edit the DOM using JQuery to display text
  function displayText(text) {
    $("#hint").text(text);
  }

  // when the button is click update status and notify sciler
  $("#button").on("click", function() {
    display.button = true;
    display.statusChanged();
  });

  // get config file from server
  $.get("/display_config.json", function(config) {
    display = new Display(JSON.parse(config)); // create new Display object
    // connect
    display.start(() => {
      console.log("connected"); // when connected, do something
    });
  });
});
