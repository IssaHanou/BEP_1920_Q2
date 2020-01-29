import {Component, OnInit, ViewChild} from "@angular/core";
import {MatSort, MatTableDataSource} from "@angular/material";
import {Event} from "./event";
import {AppComponent} from "../../app.component";

@Component({
  selector: "app-event",
  templateUrl: "./event.component.html",
  styleUrls: ["./event.component.css"]
})
export class EventComponent implements OnInit {

  /**
   * The keys used by the table to retrieve data from the DataSource
   */
  eventColumns: string[] = ["id", "status", "description", "done"];

  /**
   * Control the sorting of the table.
   */
  @ViewChild("EventTableSort", { static: true }) sort: MatSort;

  constructor(private app: AppComponent) { }

  ngOnInit() {
  }

  /**
   * Returns list of Event objects with their current status.
   * Return in the form of map table data source, with sorting enabled.
   */
  public getEventStatus(): MatTableDataSource<Event> {
    const puzzles: Event[] = [];
    for (const puzzle of this.app.puzzleList.all.values()) {
      if (!puzzle.isPuzzle) {
        puzzles.push(puzzle);
      }
    }
    puzzles.sort((a: Event, b: Event) => a.id.localeCompare(b.id));

    const dataSource = new MatTableDataSource<Event>(puzzles);
    dataSource.sort = this.sort;
    return dataSource;
  }
}
