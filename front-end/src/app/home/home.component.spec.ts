import { async, ComponentFixture, TestBed } from "@angular/core/testing";

import { HomeComponent } from "./home.component";
import { BrowserModule } from "@angular/platform-browser";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MqttModule, MqttService } from "ngx-mqtt";
import { APP_ROUTES, MQTT_SERVICE_OPTIONS } from "../app.module";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import {
  MatButtonModule,
  MatFormFieldModule,
  MatIconModule,
  MatInputModule,
  MatListModule,
  MatSelectModule,
  MatTooltipModule,
  MatSidenavModule,
  MatSnackBar,
  MatSnackBarModule,
  MatSortModule,
  MatCheckboxModule,
  MatExpansionModule,
  MatTableModule,
  MatToolbarModule
} from "@angular/material";
import { CdkTableModule } from "@angular/cdk/table";
import { RouterModule } from "@angular/router";
import { AppComponent } from "../app.component";
import { CameraComponent } from "../camera/camera.component";
import { ConfigComponent } from "../config/config.component";
import { HintComponent } from "../components/hint/hint.component";
import { DeviceComponent } from "../components/device/device.component";
import { ManageComponent } from "../components/manage/manage.component";
import { PuzzleComponent } from "../components/puzzle/puzzle.component";
import { TimerComponent } from "../components/timer/timer.component";
import { Overlay } from "@angular/cdk/overlay";
import { EventComponent } from "../components/event/event.component";

describe("HomeComponent", () => {
  let component: HomeComponent;
  let fixture: ComponentFixture<HomeComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        BrowserModule,
        FormsModule,
        MqttModule.forRoot(MQTT_SERVICE_OPTIONS),
        BrowserAnimationsModule,
        MatSnackBarModule,
        MatTableModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        MatTooltipModule,
        MatSortModule,
        MatSidenavModule,
        MatToolbarModule,
        MatIconModule,
        MatCheckboxModule,
        MatSelectModule,
        MatExpansionModule,
        MatListModule,
        CdkTableModule,
        RouterModule.forRoot(APP_ROUTES),
        ReactiveFormsModule
      ],
      declarations: [
        AppComponent,
        CameraComponent,
        HomeComponent,
        ConfigComponent,
        HintComponent,
        DeviceComponent,
        ManageComponent,
        EventComponent,
        PuzzleComponent,
        TimerComponent
      ],
      providers: [MqttService, MatSnackBar, AppComponent, Overlay]
    }).compileComponents();
  }));

  it("should create", () => {
    fixture = TestBed.createComponent(HomeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });
});
