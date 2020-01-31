import { async, ComponentFixture, TestBed } from "@angular/core/testing";

import { PuzzleComponent } from "./puzzle.component";
import {
  MatTableModule,
  MatTooltipModule,
  MatIconModule,
  MatFormFieldModule,
  MatSelectModule,
  MatInputModule,
  MatOptionModule,
  MatExpansionModule,
  MatDividerModule
} from "@angular/material";
import { MqttModule, MqttService } from "ngx-mqtt";
import { AppComponent } from "../../app.component";
import { Overlay } from "@angular/cdk/overlay";
import { MQTT_SERVICE_OPTIONS } from "../../app.module";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { MatSnackBar } from "@angular/material/snack-bar";
import { FormsModule } from "@angular/forms";

describe("PuzzleComponent", () => {
  let component: PuzzleComponent;
  let fixture: ComponentFixture<PuzzleComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        FormsModule,
        MqttModule.forRoot(MQTT_SERVICE_OPTIONS),
        MatTableModule,
        MatTooltipModule,
        MatFormFieldModule,
        MatSelectModule,
        MatInputModule,
        MatOptionModule,
        MatIconModule,
        MatDividerModule,
        MatExpansionModule,
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
