export class Timer {
  id: string;
  duration: number;
  state: string;

  constructor(jsonData) {
    this.id = jsonData.id;
    this.duration = jsonData.status;
    this.state = jsonData.state;
  }

  getState() {
    return this.state;
  }

  getTimeLeft() {
    return this.duration;
  }

  tick() {
    this.duration = this.duration - 1000;
  }

  update(dur, sta) {
    this.duration = dur;
    this.state = sta;
  }
}
