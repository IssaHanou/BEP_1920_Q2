import {Component, Injectable, OnDestroy, OnInit} from '@angular/core';
import {IMqttMessage, MqttService} from 'ngx-mqtt';

@Component({
  selector: 'app-hint',
  templateUrl: './hint.component.html',
  styleUrls: ['./hint.component.css', '../../../assets/css/main.css']
})
export class HintComponent implements OnDestroy {
  hint: string;

  private subscription;
  public message: string;

  constructor(private mqttService: MqttService) {
    this.subscription = this.mqttService.observe('test').subscribe((message: IMqttMessage) => {
      this.message = message.payload.toString();
    });
  }

  public unsafePublish(topic: string, message: string): void {
      this.mqttService.unsafePublish(topic, message, {qos: 1, retain: true});
    }

  public ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  onSubmit() {
    this.unsafePublish('test', this.hint);
    this.hint = '';
  }
}
