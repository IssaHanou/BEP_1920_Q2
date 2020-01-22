import { Component, OnInit, ViewChild } from "@angular/core";
import { MatSort, MatTableDataSource } from "@angular/material";
import { AppComponent } from "../../app.component";
import { Puzzle } from "./puzzle";

/**
 * The puzzle component controls the puzzles tables and is shown in the "Puzzels" box on the home page.
 */
@Component({
  selector: "app-puzzle",
  templateUrl: "./puzzle.component.html",
  styleUrls: ["./puzzle.component.css", "../../../assets/css/main.css"]
})
export class PuzzleComponent implements OnInit {
  /**
   * The keys used by the table to retrieve data from the DataSource
   */
  puzzleColumns: string[] = ["id", "status", "description", "done"];

  /**
   * Control the sorting of the table.
   */
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
   * When button in the table is pressed, manually override the finished status of rule in back-end.
   */
  finishRule(ruleId: string) {
    this.app.sendInstruction([{ instruction: "finish rule", rule: ruleId }]);
  }
}
