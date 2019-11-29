import { async, ComponentFixture, TestBed } from "@angular/core/testing";
import { HintComponent } from "./hint.component";
import { FormsModule } from "@angular/forms";
import { MqttModule, MqttService } from "ngx-mqtt";
import { MQTT_SERVICE_OPTIONS } from "../../app.module";

describe("HintComponent", () => {
  let component: HintComponent;
  let fixture: ComponentFixture<HintComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [FormsModule, MqttModule.forRoot(MQTT_SERVICE_OPTIONS)],
      declarations: [HintComponent],
      providers: [MqttService]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HintComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should render header", () => {
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector("h3").textContent).toContain("Hints");
  });
});
