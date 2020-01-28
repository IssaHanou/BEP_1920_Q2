import { Events } from "./events";

describe("Puzzles", () => {
  let puzzles: Events;
  let jsonData: JSON;

  beforeEach(() => {
    jsonData = JSON.parse(`[{
          "id": "Door open",
          "status": true,
           "description": "The door opens"
        }]
    `);
    puzzles = new Events();
  });

  it("should create", () => {
    expect(puzzles).toBeTruthy();
  });

  it("should set status", () => {
    expect(puzzles.all.size).toBe(0);
    puzzles.updatePuzzles(jsonData);
    expect(puzzles.all.size).toBe(1);
    puzzles.updatePuzzles(jsonData);
    expect(puzzles.all.get("Door open").status).toBe(true);
  });

  it("should add puzzle", () => {
    const newMap = new Map<string, string>();
    expect(puzzles.all.size).toBe(0);
    puzzles.updatePuzzles({
      "id": "my rule",
      "description": "this my rule",
      "status": "true",
      "puzzle": false
    });
    expect(puzzles.all.size).toBe(1);
    expect(puzzles.all.get("my rule").status).toBe(false);
  });
});
