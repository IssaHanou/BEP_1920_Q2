import {Component, OnInit} from '@angular/core';
import * as jsscc from 'node_modules/js-scc/scc.js';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit{
  title = 'display';
  hint = '';

  ngOnInit(): void {
    jsscc.test()
  }
}

