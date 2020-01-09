const http = require("http");
const fs = require("fs");

const hostname = "localhost";
const port = "3000";

const server = http.createServer((req, res) => {
  res.statusCode = 200;
  res.setHeader("Content-Type", "text/html");
  let file;
  if (req.url === "/") {
    file = "index.html";
  } else {
    file = req.url.substr(1);
  }

  fs.readFile(file, null, function(error, data) {
    if (error) {
      res.writeHead(404);
      res.write("Whoops! File not found!");
      console.log("404:" + req.url + " not found!");
    } else {
      res.write(data);
    }
    res.end();
  });
});

server.listen(port, hostname, () => {
  console.log(`Server running at http://${hostname}:${port}/`);
});
