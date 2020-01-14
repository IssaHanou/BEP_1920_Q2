import { async, ComponentFixture, TestBed } from "@angular/core/testing";
import { CameraComponent } from "./camera.component";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import { MqttModule, MqttService } from "ngx-mqtt";
import { MQTT_SERVICE_OPTIONS } from "../app.module";
import {
  MatFormFieldModule,
  MatSelectModule,
  MatCheckboxModule,
  MatSnackBar,
  MatIconModule,
  MatSnackBarModule,
  MatExpansionModule
} from "@angular/material";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { AppComponent } from "../app.component";
import { Overlay } from "@angular/cdk/overlay";

describe("CameraComponent", () => {
  let component: CameraComponent;
  let fixture: ComponentFixture<CameraComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        FormsModule,
        MqttModule.forRoot(MQTT_SERVICE_OPTIONS),
        MatFormFieldModule,
        MatSelectModule,
        MatExpansionModule,
        MatIconModule,
        MatCheckboxModule,
        MatSnackBarModule,
        BrowserAnimationsModule,
        ReactiveFormsModule
      ],
      declarations: [CameraComponent],
      providers: [MqttService, AppComponent, MatSnackBar, Overlay]
    }).compileComponents();
  }));

  it("should create", () => {
    fixture = TestBed.createComponent(CameraComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });
});
