import { async, ComponentFixture, TestBed } from "@angular/core/testing";
import { DeviceComponent } from "./device.component";

describe("TimerComponent", () => {
  let component: DeviceComponent;
  let fixture: ComponentFixture<DeviceComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [DeviceComponent]
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
    expect(tableHeaders.item(2).textContent).toContain("Status");
  });
});
