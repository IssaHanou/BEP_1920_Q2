$(document).ready(function() {
  const Device = require("js-scc");
  let display;

  class Display extends Device {
      constructor(config) {
          super(config, function (date, level, message) {
              const formatDate = function (date) {
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
          });
          this.hint = "";
      }

      getStatus() {
          return {
              "hint": this.hint
          }
      }

      performInstruction(action) {
          switch (action.instruction) {
              case "hint": {
                  this.hint = action.value;
                  displayText(this.hint);
                  this.statusChanged();
                  break;
              }
              default: {
                  return false;
              }
          }
          return true;
      }

      test() {
          this.hint = "test";
          displayText(this.hint);
      }

      reset() {
          this.hint = "";
          displayText(this.hint);
          this.statusChanged();
      }
  }

  function displayText(text) {
      $("#hint").text(text);
  }


  // get config file from server
  $.get("/display_config.json", function(config) {
      display = new Display(JSON.parse(config));
      display.start();
  });
});
