import { Timer } from "./timer";
import { TimerComponent } from "./timer.component";

describe("TimerComponent", () => {
  let timer: Timer;

  beforeEach(() => {
    const jsonData = {
      id: "timer1",
      status: 10000,
      state: "stateIdle"
    };
    timer = new Timer(jsonData);
  });

  it("should create", () => {
    expect(timer).toBeTruthy();
    expect(timer.getState()).toBe("stateIdle");
    expect(timer.getTimeLeft()).toBe(10000);
  });

  it("should set state", () => {
    expect(timer.getState()).toBe("stateIdle");
    timer.update(10000, "stateActive");
    expect(timer.getState()).toBe("stateActive");
  });

  it("should set duration", () => {
    expect(timer.getTimeLeft()).toBe(10000);
    timer.update(2000, "stateIdle");
    expect(timer.getTimeLeft()).toBe(2000);
  });

  it("should tick with 1000 ms at the time", () => {
    timer.update(10000, "stateActive");
    expect(timer.getTimeLeft()).toBe(10000);
    timer.tick();
    expect(timer.getTimeLeft()).toBe(9000);
  });
});
