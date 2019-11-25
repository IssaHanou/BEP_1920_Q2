import { async, ComponentFixture, TestBed } from "@angular/core/testing";

import { TestComponent } from "./test.component";
import { MqttModule, MqttService } from "ngx-mqtt";
import { FormsModule } from "@angular/forms";
import { MQTT_SERVICE_OPTIONS } from "../../app.module";

describe("TestComponent", () => {
  let component: TestComponent;
  let fixture: ComponentFixture<TestComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [FormsModule, MqttModule.forRoot(MQTT_SERVICE_OPTIONS)],
      declarations: [TestComponent],
      providers: [MqttService]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TestComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
