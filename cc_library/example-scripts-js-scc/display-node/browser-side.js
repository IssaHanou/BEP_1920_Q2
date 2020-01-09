$(document).ready(function() {
  const Scclib = require("js-scc");

  console.log("Hello World!");
  $.get("/display_config.json", function(config) {
      let scc = new Scclib(JSON.parse(config), 4, function (date, level, message) {
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
          ); // call own logger
      });
  });
});
