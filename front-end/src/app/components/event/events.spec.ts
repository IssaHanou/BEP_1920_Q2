import { Events } from "./events";

describe("Puzzles", () => {
  let puzzles: Events;
  let jsonData: JSON;

  beforeEach(() => {
    jsonData = JSON.parse(`[{
          "id": "Door open",
          "status": true,
           "description": "The door opens",
           "puzzle": false
        }]
    `);
    puzzles = new Events();
  });

  it("should create", () => {
    expect(puzzles).toBeTruthy();
  });

  it("should add puzzle and set status", () => {
    puzzles.updatePuzzles(jsonData);
    expect(puzzles.all.size).toBe(1);
    puzzles.updatePuzzles(jsonData);
    expect(puzzles.all.get("Door open").status).toBe(true);
  });
});
