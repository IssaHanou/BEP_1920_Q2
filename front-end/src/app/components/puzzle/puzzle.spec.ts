import { Puzzle } from "./puzzle";

describe("Puzzle", () => {
  let puzzle: Puzzle;

  beforeEach(() => {
    puzzle = new Puzzle("Door open", "this is my rule");
  });

  it("should create", () => {
    expect(puzzle).toBeTruthy();
  });

  it("should set status", () => {
    expect(puzzle.status).toBe(false);
    puzzle.updateStatus(true);
    expect(puzzle.status).toBe(true);
  });
});
