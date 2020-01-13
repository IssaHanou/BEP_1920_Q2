import { Component, OnInit, ViewChild } from "@angular/core";
import { MatSort, MatTableDataSource } from "@angular/material";
import { AppComponent } from "../../app.component";
import { Puzzle } from "./puzzle";

@Component({
  selector: "app-puzzle",
  templateUrl: "./puzzle.component.html",
  styleUrls: ["./puzzle.component.css", "../../../assets/css/main.css"]
})
export class PuzzleComponent implements OnInit {
  puzzleColumns: string[] = ["id", "status", "description", "done"];

  @ViewChild("PuzzleTableSort", { static: true }) sort: MatSort;

  constructor(private app: AppComponent) {}

  ngOnInit() {}

  /**
   * Returns list of Puzzle objects with their current status.
   * Return in the form of map table data source, with sorting enabled.
   */
  public getPuzzleStatus(): MatTableDataSource<Puzzle> {
    const puzzles: Puzzle[] = [];
    for (const puzzle of this.app.puzzleList.all.values()) {
      puzzles.push(puzzle);
    }
    puzzles.sort((a: Puzzle, b: Puzzle) => a.id.localeCompare(b.id));

    const dataSource = new MatTableDataSource<Puzzle>(puzzles);
    dataSource.sort = this.sort;
    return dataSource;
  }

  /**
   * When button is pressed, manually override the finished status of rule in back-end.
   */
  finishRule(ruleId: string) {
    this.app.sendInstruction([{ instruction: "finish rule", rule: ruleId }]);
  }
}
