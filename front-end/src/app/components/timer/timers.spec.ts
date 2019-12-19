import { Timers } from "./timers";

describe("TimerComponent", () => {
  let timers: Timers;

  beforeEach(() => {
    timers = new Timers();
  });

  it("should create", () => {
    expect(timers).toBeTruthy();
  });

  it("should be able to add new timers", () => {
    expect(timers.getTimer("timer1")).toBeNull();
    const jsonData = {
      id: "timer1",
      status: 10000,
      state: "stateIdle"
    };
    timers.setTimer(jsonData);
    expect(timers.getTimer("timer1")).toBeTruthy();
  });

  it("should be able to update existing timers", () => {
    expect(timers.getTimer("timer1")).toBeNull();
    const jsonData = {
      id: "timer1",
      status: 10000,
      state: "stateIdle"
    };
    timers.setTimer(jsonData);
    expect(timers.getTimer("timer1").getState()).toBe("stateIdle");
    const jsonData2 = {
      id: "timer1",
      status: 200,
      state: "stateActive"
    };
    timers.setTimer(jsonData2);
    expect(timers.getTimer("timer1").getState()).toBe("stateActive");
  });
});
