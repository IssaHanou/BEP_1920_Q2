import { TestBed, async } from "@angular/core/testing";
import { AppComponent } from "./app.component";
import { HintComponent } from "./components/hint/hint.component";
import { FormsModule } from "@angular/forms";
import { MqttModule, MqttService } from "ngx-mqtt";
import { MQTT_SERVICE_OPTIONS } from "./app.module";
import { DeviceComponent } from "./components/device/device.component";
import { ManageComponent } from "./components/manage/manage.component";
import { PuzzleComponent } from "./components/puzzle/puzzle.component";
import { TimerComponent } from "./components/timer/timer.component";
import {MatSnackBarModule} from "@angular/material/snack-bar";

describe("AppComponent", () => {
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [FormsModule,
        MqttModule.forRoot(MQTT_SERVICE_OPTIONS),
        MatSnackBarModule],
      declarations: [
        AppComponent,
        HintComponent,
        DeviceComponent,
        ManageComponent,
        PuzzleComponent,
        TimerComponent
      ],
      providers: [MqttService]
    }).compileComponents();
  }));

  it("should create the app", () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app).toBeTruthy();
  });

  it("should have as title 'S.C.I.L.E.R'", () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app.title).toEqual("S.C.I.L.E.R");
  });

  it(`should have as subtitle 'Super awesome escape'`, () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app.nameOfRoom).toEqual("Super awesome escape");
  });

  it("should render title", () => {
    const fixture = TestBed.createComponent(AppComponent);
    fixture.detectChanges();
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector("h1").textContent).toContain("S.C.I.L.E.R");
  });

  it("should render subtitle", () => {
    const fixture = TestBed.createComponent(AppComponent);
    fixture.detectChanges();
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector("h2").textContent).toContain(
      "Super awesome escape"
    );
  });
});
