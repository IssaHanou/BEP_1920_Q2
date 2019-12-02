import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ManageComponent } from "./manage.component";
import { MqttModule, MqttService } from "ngx-mqtt";
import { MQTT_SERVICE_OPTIONS } from "../../app.module";

describe("ManageComponent", () => {
  let component: ManageComponent;
  let fixture: ComponentFixture<ManageComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [MqttModule.forRoot(MQTT_SERVICE_OPTIONS)],
      declarations: [ManageComponent],
      providers: [MqttService]
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ManageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should render header", () => {
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector("h3").textContent).toContain("Acties");
  });

  it("should have all buttons present", () => {
    const compiled = fixture.debugElement.nativeElement;
    expect(
      compiled.getElementsByClassName("start").item(0).textContent
    ).toContain("Start");
    expect(
      compiled.getElementsByClassName("test").item(0).textContent
    ).toContain("Test");
    expect(
      compiled.getElementsByClassName("stop").item(0).textContent
    ).toContain("Stop");
    expect(
      compiled.getElementsByClassName("reset").item(0).textContent
    ).toContain("Reset");
  });

  it("should send test message on test button click", () => {
    spyOn(component, "onClickTestButton");
    const testButton = fixture.debugElement.nativeElement
      .getElementsByClassName("test")
      .item(0);
    testButton.click();
    fixture.detectChanges();

    expect(component.onClickTestButton).toHaveBeenCalled();
  });
});
