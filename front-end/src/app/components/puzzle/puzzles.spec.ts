import { Puzzles } from "./puzzles";
import { PuzzleComponent } from "./puzzle.component";

describe("PuzzleComponent", () => {
  let puzzles: Puzzles;
  let jsonData: JSON;

  beforeEach(() => {
    jsonData = JSON.parse(`[{
          "id": "Door open",
          "status": true,
           "description": "The door opens"
        }]
    `);
    puzzles = new Puzzles();
  });

  it("should create", () => {
    expect(puzzles).toBeTruthy();
  });

  it("should set status", () => {
    expect(puzzles.all.size).toBe(0);
    puzzles.updatePuzzles(jsonData);
    expect(puzzles.all.size).toBe(1);
    expect(puzzles.all.get("Door open").status).toBe(true);
  });
});
