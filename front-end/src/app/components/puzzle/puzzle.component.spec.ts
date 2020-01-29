import { async, ComponentFixture, TestBed } from "@angular/core/testing";

import { PuzzleComponent } from "./puzzle.component";
import {
  MatPaginatorModule,
  MatTooltipModule,
  MatSnackBar,
  MatSortModule,
  MatFormFieldModule,
  MatInputModule,
  MatOptionModule,
  MatSelectModule,
  MatExpansionModule,
  MatDividerModule,
  MatIconModule,
  MatTableModule
} from "@angular/material";
import { MqttModule, MqttService } from "ngx-mqtt";
import { AppComponent } from "../../app.component";
import { Overlay } from "@angular/cdk/overlay";
import { MQTT_SERVICE_OPTIONS } from "../../app.module";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";

describe("PuzzleComponent", () => {
  let component: PuzzleComponent;
  let fixture: ComponentFixture<PuzzleComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        MqttModule.forRoot(MQTT_SERVICE_OPTIONS),
        MatTableModule,
        MatTooltipModule,
        MatSortModule,
        MatFormFieldModule,
        MatInputModule,
        MatOptionModule,
        MatSelectModule,
        MatIconModule,
        MatDividerModule,
        MatExpansionModule,
        MatPaginatorModule,
        BrowserAnimationsModule
      ],
      declarations: [PuzzleComponent],
      providers: [MqttService, AppComponent, MatSnackBar, Overlay]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PuzzleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should render header", () => {
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector("h3").textContent).toContain("Puzzels");
  });
});
