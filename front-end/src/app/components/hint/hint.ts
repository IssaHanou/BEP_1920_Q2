export class Hint {
  puzzle: string;
  hints: string[];

  constructor(puzzle: string, hints: string[]) {
    this.puzzle = puzzle;
    this.hints = hints;
  }
}
