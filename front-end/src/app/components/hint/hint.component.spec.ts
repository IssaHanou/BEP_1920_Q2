import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { HintComponent } from './hint.component';
import {FormsModule} from '@angular/forms';
import {MqttService} from 'ngx-mqtt';

describe('HintComponent', () => {
  let component: HintComponent;
  let fixture: ComponentFixture<HintComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [ FormsModule ],
      declarations: [ HintComponent ],
      providers: [ MqttService ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HintComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
