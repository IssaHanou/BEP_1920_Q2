import { BrowserModule, HAMMER_LOADER } from "@angular/platform-browser";
import { NgModule } from "@angular/core";
import { AppComponent } from "./app.component";
import { HintComponent } from "./components/hint/hint.component";
import { DeviceComponent } from "./components/device/device.component";
import { TimerComponent } from "./components/timer/timer.component";
import { ManageComponent } from "./components/manage/manage.component";
import { PuzzleComponent } from "./components/puzzle/puzzle.component";
import { FormsModule } from "@angular/forms";

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
  MatPaginatorModule,
  MatSelectModule,
  MatSortModule,
  MatTableModule
} from "@angular/material";
import { CdkTableModule } from "@angular/cdk/table";
import { JsonConvert } from "json2typescript";
import { Message } from "./message";

export const MQTT_SERVICE_OPTIONS: IMqttServiceOptions = {
  hostname: "192.168.178.82",
  port: 8083,
  clientId: "front-end",
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
  keepalive: 10
};

@NgModule({
  declarations: [
    AppComponent,
    HintComponent,
    DeviceComponent,
    TimerComponent,
    ManageComponent,
    PuzzleComponent
  ],
  exports: [
    AppComponent,
    HintComponent,
    DeviceComponent,
    TimerComponent,
    ManageComponent,
    PuzzleComponent,
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
    MatFormFieldModule,
    MatInputModule,
    MatPaginatorModule,
    MatSortModule,
    MatSelectModule,
    CdkTableModule
  ],
  providers: [
    MqttService,
    MatSnackBar,
    Overlay,
    DeviceComponent,
    {
      provide: HAMMER_LOADER,
      useValue: () => new Promise(() => {})
    }, // prevents warning in console
    { provide: MATERIAL_SANITY_CHECKS, useValue: false } // prevents warning in console
  ],
  bootstrap: [AppComponent],
  entryComponents: [MatSnackBarContainer]
})
export class AppModule {}
