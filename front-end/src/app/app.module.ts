import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";
import { AppComponent } from "./app.component";
import { HintComponent } from "./components/hint/hint.component";
import { DeviceComponent } from "./components/device/device.component";
import { TimerComponent } from "./components/timer/timer.component";
import { ManageComponent } from "./components/manage/manage.component";
import { PuzzleComponent } from "./components/puzzle/puzzle.component";
import { FormsModule } from "@angular/forms";

import { MqttModule, MqttService, IMqttServiceOptions } from "ngx-mqtt";
import {MatSnackBar, MatSnackBarContainer, MatSnackBarModule} from "@angular/material/snack-bar";
import {Overlay} from "@angular/cdk/overlay";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";

export const MQTT_SERVICE_OPTIONS: IMqttServiceOptions = {
  hostname: "192.168.178.82",
  port: 8083,
  clientId: "front-end"
};

@NgModule({
  declarations: [
    AppComponent,
    HintComponent,
    DeviceComponent,
    TimerComponent,
    ManageComponent,
    PuzzleComponent,
  ],
  exports: [
    AppComponent,
    HintComponent,
    DeviceComponent,
    TimerComponent,
    ManageComponent,
    PuzzleComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    MqttModule.forRoot(MQTT_SERVICE_OPTIONS),
    BrowserAnimationsModule,
    MatSnackBarModule
  ],
  providers: [MqttService, MatSnackBar, Overlay],
  bootstrap: [AppComponent],
  entryComponents: [MatSnackBarContainer]
})
export class AppModule {}
