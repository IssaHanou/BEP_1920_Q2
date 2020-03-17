import { Component } from '@angular/core';
import {Slides} from "./slides";
import * as data from "../assets/slides_config.json";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = "slides";
  config;
  slides;

  constructor() {
    this.config = (data as any).default;
    this.slides = new Slides(this.config);
  }

  ngOnInit(): void {
    this.slides.start(() => {
      console.log("CONNECTED");
    });
  }

  getSrc() {
    if (this.slides.allSlides === undefined) {
      return "../assets/images/Opening.jpg";
    } else {
      return "../assets/images/" + this.slides.allSlides[this.slides.current] + ".jpg";
    }
  }
}
