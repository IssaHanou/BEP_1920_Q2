import { async, ComponentFixture, TestBed } from "@angular/core/testing";

import { EventComponent } from "./event.component";
import {
  MatPaginatorModule,
  MatTooltipModule,
  MatSnackBar,
  MatIconModule,
  MatSortModule,
  MatExpansionModule,
  MatTableModule
} from "@angular/material";
import { MqttModule, MqttService } from "ngx-mqtt";
import { AppComponent } from "../../app.component";
import { Overlay } from "@angular/cdk/overlay";
import { MQTT_SERVICE_OPTIONS } from "../../app.module";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";

describe("EventComponent", () => {
  let component: EventComponent;
  let fixture: ComponentFixture<EventComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [EventComponent],
      imports: [
        MqttModule.forRoot(MQTT_SERVICE_OPTIONS),
        MatTableModule,
        MatTooltipModule,
        MatIconModule,
        MatExpansionModule,
        MatSortModule,
        MatPaginatorModule,
        BrowserAnimationsModule
      ],
      providers: [MqttService, AppComponent, MatSnackBar, Overlay]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EventComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should render header", () => {
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector("h3").textContent).toContain(
      "Gebeurtenissen"
    );
  });
});
