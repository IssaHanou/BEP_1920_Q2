import { Component, OnInit, ViewChild } from "@angular/core";
import { MatPaginator, MatSort, MatTableDataSource } from "@angular/material";
import { AppComponent } from "../../app.component";

export interface Puzzle {
  id: string
  status: boolean
}

@Component({
  selector: "app-puzzle",
  templateUrl: "./puzzle.component.html",
  styleUrls: ["./puzzle.component.css", "../../../assets/css/main.css"]
})
export class PuzzleComponent implements OnInit {
  puzzleColumns: string[] = ["id", "status"];

  @ViewChild("PuzzleTableSort", {static: true}) sort: MatSort;

  constructor(private app: AppComponent) {}

  ngOnInit() {
  }

  /**
   * Returns list of Device object with their current status and connection.
   * Return in the form of map table data source, with sorting enabled.
   */
  public getPuzzleStatus(): MatTableDataSource<Puzzle> {
    const puzzles: Puzzle[] = [
      {id: "Telefoon puzzle", status: true},
      {id: "controle board", status: false},
      {id: "yet another", status: false}
    ];
    puzzles.sort((a: Puzzle, b: Puzzle) => a.id.localeCompare(b.id));

    const dataSource = new MatTableDataSource<Puzzle>(puzzles);
    dataSource.sort = this.sort;
    return dataSource;
  }
}
