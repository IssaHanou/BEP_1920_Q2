import { Component, OnInit } from "@angular/core";

export interface Puzzle {
  puzzle: string
  status: boolean
}

@Component({
  selector: "app-puzzle",
  templateUrl: "./puzzle.component.html",
  styleUrls: ["./puzzle.component.css", "../../../assets/css/main.css"]
})
export class PuzzleComponent implements OnInit {
  puzzleColumns: string[] = ['puzzel', 'status'];
  puzzleData: Puzzle[] = [
    {puzzle: "Telefoon puzzle", status: true},
    {puzzle: "controle board", status: false},
    {puzzle: "yet another", status: false}
  ];

  constructor() {}

  ngOnInit() {}
}
