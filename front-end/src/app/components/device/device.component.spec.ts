import { async, ComponentFixture, TestBed } from "@angular/core/testing";
import { DeviceComponent } from "./device.component";
import { MqttModule, MqttService } from "ngx-mqtt";
import { MQTT_SERVICE_OPTIONS } from "../../app.module";
import { FormsModule } from "@angular/forms";
import { AppComponent } from "../../app.component";
import { MatSnackBar } from "@angular/material/snack-bar";
import { Overlay } from "@angular/cdk/overlay";
import {
  MatFormFieldModule,
  MatInputModule,
  MatPaginatorModule,
  MatTooltipModule,
  MatIconModule,
  MatSortModule,
  MatExpansionModule,
  MatTableModule
} from "@angular/material";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";

describe("DeviceComponent", () => {
  let component: DeviceComponent;
  let fixture: ComponentFixture<DeviceComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        FormsModule,
        MqttModule.forRoot(MQTT_SERVICE_OPTIONS),
        MatTableModule,
        MatSortModule,
        MatPaginatorModule,
        MatFormFieldModule,
        MatInputModule,
        MatTooltipModule,
        MatIconModule,
        MatExpansionModule,
        BrowserAnimationsModule
      ],
      declarations: [DeviceComponent],
      providers: [MqttService, AppComponent, MatSnackBar, Overlay]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DeviceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should render header", () => {
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector("h3").textContent).toContain("Apparaten");
  });

  it("table should have correct headers", () => {
    const compiled = fixture.debugElement.nativeElement;
    const tableHeaders = compiled.querySelectorAll("th");
    expect(tableHeaders.item(0).textContent).toContain("Apparaat");
    expect(tableHeaders.item(1).textContent).toContain("Connectie");
    expect(tableHeaders.item(2).textContent).toContain("Laat status zien");
    expect(tableHeaders.item(3).textContent).toContain("Test");
  });
});
