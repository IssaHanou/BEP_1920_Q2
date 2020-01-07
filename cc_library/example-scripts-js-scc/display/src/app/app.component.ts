import { Component, OnInit } from "@angular/core";
import { SccLib } from "node_modules/js-scc/scc.js";
import { HttpClient } from "@angular/common/http";


@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.css"]
})
export class AppComponent implements OnInit {
  title = "display";
  hint = "";
  scc;

  constructor(private http: HttpClient) {}

  ngOnInit(): void {
    this.http
      .get("assets/display_config.json")
      .toPromise()
      .then((response: any) => {
        const config = response;
        this.scc = new SccLib(config, 4);
        this.scc.test()
      });

  }
}
