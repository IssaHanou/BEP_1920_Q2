import { Component } from '@angular/core';
import { Navigation } from "./navigation";
import * as data from "../assets/navigate_config.json";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = "navigate";
  config;
  navigation;

  constructor() {
    this.config = (data as any).default;
    this.navigation = new Navigation(this.config);
  }

  ngOnInit(): void {
    this.navigation.start(() => {
      console.log("CONNECTED");
    });
  }
}
