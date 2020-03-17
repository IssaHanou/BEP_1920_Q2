import { Component } from '@angular/core';
import {Navigation} from "../../../navigate/src/app/navigation";
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
    this.slides = new Navigation(this.config);
  }

  ngOnInit(): void {
    this.slides.start(() => {
      console.log("CONNECTED");
    });
  }

  getSrc() {
    if (this.slides.allSlides === undefined) {
      return "../assets/images/opening.JPG";
    } else {
      return "../assets/images/" + this.slides.allSlides[this.slides.current] + ".JPG";
    }
  }
}
