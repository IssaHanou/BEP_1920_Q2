import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { EventComponent } from './event.component';
import {MqttModule, MqttService} from "ngx-mqtt";
import {MQTT_SERVICE_OPTIONS} from "../../app.module";
import {MatPaginatorModule, MatSnackBar, MatSortModule, MatTableModule} from "@angular/material";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {PuzzleComponent} from "../puzzle/puzzle.component";
import {AppComponent} from "../../app.component";
import {Overlay} from "@angular/cdk/overlay";

describe('EventComponent', () => {
  let component: EventComponent;
  let fixture: ComponentFixture<EventComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ EventComponent ],
      imports: [
        MqttModule.forRoot(MQTT_SERVICE_OPTIONS),
        MatTableModule,
        MatSortModule,
        MatPaginatorModule,
        BrowserAnimationsModule],
      providers: [MqttService, AppComponent, MatSnackBar, Overlay]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EventComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it("should render header", () => {
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector("h3").textContent).toContain("Gebeurtenissen");
  });

  it("table should have correct headers", () => {
    const compiled = fixture.debugElement.nativeElement;
    const tableHeaders = compiled.querySelectorAll("th");
    expect(tableHeaders.item(0).textContent).toContain("Gebeurtenis");
    expect(tableHeaders.item(1).textContent).toContain("Uitgevoerd");
    expect(tableHeaders.item(2).textContent).toContain("Beschrijving");
    expect(tableHeaders.item(3).textContent).toContain("Handmatig afmaken");
  });
});
