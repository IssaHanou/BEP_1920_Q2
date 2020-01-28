import { Event } from "./event";

describe("Puzzle", () => {
  let puzzle: Event;

  beforeEach(() => {
    puzzle = new Event("Door open", "this is my rule");
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
