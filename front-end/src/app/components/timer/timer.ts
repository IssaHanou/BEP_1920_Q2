export class Timer {
  id: string;
  duration: number;
  state: string;

  constructor(jsonData) {
    this.id = jsonData.id;
    this.duration = jsonData.duration;
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

/**
 * Format the time in milliseconds to a string in the format hh:mm:ss.
 */
export function formatMS(timeInMS) {
  const seconds = parseInt(((timeInMS / 1000) % 60).toString(), 10);
  const minutes = parseInt(((timeInMS / (1000 * 60)) % 60).toString(), 10);
  const hours = parseInt(((timeInMS / (1000 * 60 * 60)) % 24).toString(), 10);
  const h = hours < 10 ? "0" + hours : hours;
  const m = minutes < 10 ? "0" + minutes : minutes;
  const s = seconds < 10 ? "0" + seconds : seconds;

  return h + ":" + m + ":" + s;
}

/**
 * Format the time in milliseconds to a string in the format hh:mm.
 * Timezone off set is in minutes.
 */
export function formatTime(timeInMS, tzOffSet) {
  const tzInMs = 60 * 1000 * tzOffSet;
  timeInMS = timeInMS - tzInMs;
  return formatMS(timeInMS).substr(0, 5);
}
