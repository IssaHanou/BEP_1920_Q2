import { async, ComponentFixture, TestBed } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MqttModule, MqttService } from "ngx-mqtt";
import { MQTT_SERVICE_OPTIONS } from "../app.module";
import {
  MatExpansionModule,
  MatFormFieldModule,
  MatSelectModule,
  MatListModule,
  MatSnackBar,
  MatSnackBarModule
} from "@angular/material";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { AppComponent } from "../app.component";
import { ConfigComponent } from "./config.component";
import { Overlay } from "@angular/cdk/overlay";

describe("ConfigComponent", () => {
  let component: ConfigComponent;
  let fixture: ComponentFixture<ConfigComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        FormsModule,
        MqttModule.forRoot(MQTT_SERVICE_OPTIONS),
        MatFormFieldModule,
        MatSelectModule,
        MatListModule,
        MatExpansionModule,
        MatSnackBarModule,
        BrowserAnimationsModule,
        ReactiveFormsModule
      ],
      declarations: [ConfigComponent],
      providers: [MqttService, AppComponent, MatSnackBar, Overlay]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ConfigComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
