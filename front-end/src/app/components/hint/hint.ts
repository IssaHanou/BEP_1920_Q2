/**
 * Hint object, which contains a list of all hints that belong to a certain puzzle
 */
export class Hint {
  puzzle: string;
  hints: string[];

  constructor(puzzle: string, hints: string[]) {
    this.puzzle = puzzle;
    this.hints = hints;
  }
}
