import { Event } from "./event";

describe("Event", () => {
  let puzzle: Event;

  beforeEach(() => {
    puzzle = new Event("Door open", "this is my rule", false, "name", true);
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
