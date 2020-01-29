import { BrowserModule, HAMMER_LOADER } from "@angular/platform-browser";
import { NgModule } from "@angular/core";
import { AppComponent } from "./app.component";
import { HintComponent } from "./components/hint/hint.component";
import { DeviceComponent } from "./components/device/device.component";
import { TimerComponent } from "./components/timer/timer.component";
import { ManageComponent } from "./components/manage/manage.component";
import { PuzzleComponent } from "./components/puzzle/puzzle.component";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import * as config from "../assets/config.json";

import { MqttModule, MqttService, IMqttServiceOptions } from "ngx-mqtt";
import {
  MatSnackBar,
  MatSnackBarContainer,
  MatSnackBarModule
} from "@angular/material/snack-bar";
import { Overlay } from "@angular/cdk/overlay";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import {
  MatButtonModule,
  MATERIAL_SANITY_CHECKS,
  MatFormFieldModule,
  MatInputModule,
  MatIconModule,
  MatTooltipModule,
  MatCheckboxModule,
  MatSortModule,
  MatTableModule,
  MatToolbarModule,
  MatSidenavModule,
  MatListModule,
  MatSelectModule,
  MatExpansionModule
} from "@angular/material";
import { CdkTableModule } from "@angular/cdk/table";
import { JsonConvert } from "json2typescript";
import { Message } from "./message";
import { RouterModule, Routes } from "@angular/router";
import { HomeComponent } from "./home/home.component";
import { CameraComponent } from "./camera/camera.component";
import { ConfigComponent } from "./config/config.component";
import { EventComponent } from "./components/event/event.component";

/**
 * These are the parameters used for the MQTT messaging.
 * The host and port are read from ./../assets/config.json
 */
export const MQTT_SERVICE_OPTIONS: IMqttServiceOptions = {
  hostname: config.host,
  port: config.port,
  will: {
    topic: "back-end",
    payload: JSON.stringify(
      new JsonConvert().serialize(
        new Message("front-end", "connection", new Date(), {
          connection: false
        })
      )
    ),
    qos: 1,
    retain: false
  },
  keepalive: 5
};

/**
 * These are the routes used in the application: there is the home page, camera page and config page.
 */
export const APP_ROUTES: Routes = [
  { path: "", component: HomeComponent },
  { path: "camera", component: CameraComponent },
  { path: "config", component: ConfigComponent }
];

/**
 * This module runs the application and controls all the imports.
 * The application is started from the ./../index.html file, which adds contents from the application component.
 * In the ./../assets/css/main.css file most of the css general to the complete application is located.
 * Other local css can be found in each component.
 */
@NgModule({
  declarations: [
    AppComponent,
    HintComponent,
    DeviceComponent,
    TimerComponent,
    ManageComponent,
    PuzzleComponent,
    HomeComponent,
    CameraComponent,
    ConfigComponent,
    EventComponent
  ],
  exports: [
    AppComponent,
    HintComponent,
    DeviceComponent,
    TimerComponent,
    ManageComponent,
    PuzzleComponent,
    HomeComponent,
    CameraComponent,
    ConfigComponent,
    MatFormFieldModule,
    MatSortModule
  ],
  imports: [
    BrowserModule,
    FormsModule,
    MqttModule.forRoot(MQTT_SERVICE_OPTIONS),
    BrowserAnimationsModule,
    MatSnackBarModule,
    MatTableModule,
    MatButtonModule,
    MatCheckboxModule,
    MatFormFieldModule,
    MatInputModule,
    MatSortModule,
    MatTooltipModule,
    MatSidenavModule,
    MatToolbarModule,
    MatIconModule,
    MatSelectModule,
    MatExpansionModule,
    MatListModule,
    RouterModule.forRoot(APP_ROUTES),
    ReactiveFormsModule,
    CdkTableModule
  ],
  providers: [
    MqttService,
    MatSnackBar,
    Overlay,
    {
      provide: HAMMER_LOADER,
      useValue: () => new Promise(() => {})
    }, // prevents warning in console
    {
      provide: MATERIAL_SANITY_CHECKS,
      useValue: false
    } // prevents warning in console
  ],
  bootstrap: [AppComponent],
  entryComponents: [MatSnackBarContainer]
})
export class AppModule {}
