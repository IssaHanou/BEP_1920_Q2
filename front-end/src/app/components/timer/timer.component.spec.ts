import {async, ComponentFixture, TestBed} from "@angular/core/testing";

import {TimerComponent} from "./timer.component";
import {AppComponent} from "../../app.component";
import {MqttModule, MqttService} from "ngx-mqtt";
import {MatSnackBar} from "@angular/material/snack-bar";
import {Overlay} from "@angular/cdk/overlay";
import {FormsModule} from "@angular/forms";
import {MQTT_SERVICE_OPTIONS} from "../../app.module";
import {MatFormFieldModule, MatInputModule} from "@angular/material";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";

describe("TimerComponent", () => {
  let component: TimerComponent;
  let fixture: ComponentFixture<TimerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        FormsModule,
        MqttModule.forRoot(MQTT_SERVICE_OPTIONS),
        MatFormFieldModule,
        MatInputModule,
        BrowserAnimationsModule
      ],
      declarations: [TimerComponent],
      providers: [MqttService, AppComponent, MatSnackBar, Overlay]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TimerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should render header", () => {
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector("h3").textContent).toContain("Tijd");
  });
});
