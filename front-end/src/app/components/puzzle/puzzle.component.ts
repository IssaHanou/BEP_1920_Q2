import {Component, OnInit, ViewChild} from "@angular/core";
import {MatSort, MatTableDataSource} from "@angular/material";

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
  puzzleColumns: string[] = ['puzzle', 'status'];
  puzzleData = new MatTableDataSource([
    {puzzle: "Telefoon puzzle", status: true},
    {puzzle: "controle board", status: false},
    {puzzle: "yet another", status: false}
  ]);

  @ViewChild(MatSort, {static: true}) sort: MatSort;

  ngOnInit() {
    this.puzzleData.sort = this.sort
  }
}
