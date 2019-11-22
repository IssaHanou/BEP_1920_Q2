import { BrowserModule } from '@angular/platform-browser';
import {NgModule, OnDestroy} from '@angular/core';

import { AppComponent } from './app.component';
import { HintComponent } from './components/hint/hint.component';
import { FormsModule } from '@angular/forms';
import { Observable } from 'rxjs';

import {
  IMqttMessage,
  MqttModule,
  MqttService,
  IMqttServiceOptions
} from 'ngx-mqtt';

export const MQTT_SERVICE_OPTIONS: IMqttServiceOptions = {
  hostname: 'localhost',
  port: 8083,
};


@NgModule({
  declarations: [
    AppComponent,
    HintComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    MqttModule.forRoot(MQTT_SERVICE_OPTIONS)
  ],
  providers: [ MqttService ],
  bootstrap: [AppComponent]
})
export class AppModule { }
