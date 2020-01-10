import { TestBed, async } from "@angular/core/testing";
import { AppComponent } from "./app.component";
import { HintComponent } from "./components/hint/hint.component";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MqttModule, MqttService } from "ngx-mqtt";
import { APP_ROUTES, MQTT_SERVICE_OPTIONS } from "./app.module";
import { DeviceComponent } from "./components/device/device.component";
import { ManageComponent } from "./components/manage/manage.component";
import { PuzzleComponent } from "./components/puzzle/puzzle.component";
import { TimerComponent } from "./components/timer/timer.component";
import { MatSnackBar, MatSnackBarModule } from "@angular/material/snack-bar";
import { Overlay } from "@angular/cdk/overlay";
import {
  MatFormFieldModule,
  MatIconModule,
  MatInputModule,
  MatSidenavModule,
  MatTableModule,
  MatToolbarModule,
  MatListModule,
  MatButtonModule,
  MatSelectModule,
  MatExpansionModule,
  MatSortModule
} from "@angular/material";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { BrowserModule } from "@angular/platform-browser";
import { CdkTableModule } from "@angular/cdk/table";
import { RouterModule } from "@angular/router";
import { CameraComponent } from "./camera/camera.component";
import { HomeComponent } from "./home/home.component";
import { ConfigComponent } from "./config/config.component";

describe("AppComponent", () => {
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
        MatSelectModule,
        MatInputModule,
        MatSortModule,
        MatSidenavModule,
        MatToolbarModule,
        MatIconModule,
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
        PuzzleComponent,
        TimerComponent
      ],
      providers: [MqttService, MatSnackBar, Overlay]
    }).compileComponents();
  }));

  it("should create the app", () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app).toBeTruthy();
  });

  it("should have as title 'SCILER'", () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app.title).toEqual("SCILER");
  });

  it("should render title", () => {
    const fixture = TestBed.createComponent(AppComponent);
    fixture.detectChanges();
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelectorAll("p").item(0).textContent).toContain("SCILER");
  });
});
