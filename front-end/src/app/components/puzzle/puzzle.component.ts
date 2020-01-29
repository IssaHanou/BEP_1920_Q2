import { Component, OnInit, ViewChild } from "@angular/core";
import { MatSort, MatTableDataSource } from "@angular/material";
import { AppComponent } from "../../app.component";
import { Event } from "../event/event";

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
   * Hint to send for a certain puzzle.
   */
  hint: string;

  /**
   * Id of device to which a hint must be sent.
   */
  device: string;

  /**
   * Control the sorting of the table.
   */
  @ViewChild("PuzzleTableSort", {static: true}) sort: MatSort;

  constructor(private app: AppComponent) {
  }

  ngOnInit() {
  }

  /**
   * Returns list of Event objects with their current status.
   * Return in the form of map table data source, with sorting enabled.
   */
  public getPuzzleStatus(): MatTableDataSource<Event> {
    const puzzles: Event[] = [];
    for (const puzzle of this.app.puzzleList.all.values()) {
      if (puzzle.isPuzzle) {
        puzzles.push(puzzle);
      }
    }
    puzzles.sort((a: Event, b: Event) => a.id.localeCompare(b.id));

    const dataSource = new MatTableDataSource<Event>(puzzles);
    dataSource.sort = this.sort;
    return dataSource;
  }


  /**
   * Hint list used for selection of predefined hints.
   * This is generated each time from the app hint list, to ensure updated version.
   */
  getHintList(puzzle: string): string[] {
    const list = [];
    for (const obj of this.app.hintList) {
      if (obj.puzzle === puzzle) {
        for (const hint of obj.hints) {
          list.push(hint);
        }
        return list;
      }
    }
  }


  /**
   * When a hint has been chosen for a puzzle and the accompanying "Stuur" button is clicked,
   * the selected hint is sent as instruction to hint devices.
   */
  onSelectedHint() {
    if (
      this.hint !== undefined &&
      this.hint !== ""
    ) {
      this.app.sendInstruction([
        {instruction: "hint", value: this.hint}

      ]);
      this.hint = "";
    }
  }

  /**
   * Get list of devices which can show hints.
   */
  getHintDevices() {
    return [];
  }
}
