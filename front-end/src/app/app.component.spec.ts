import { TestBed, async } from '@angular/core/testing';
import { AppComponent } from './app.component';
import {HintComponent} from './components/hint/hint.component';
import {FormsModule} from '@angular/forms';
import {MqttService} from 'ngx-mqtt';

describe('AppComponent', () => {
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [ FormsModule ],
      declarations: [
        AppComponent,
        HintComponent,
      ],
      providers: [ MqttService ]
    }).compileComponents();
  }));

  it('should create the app', () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app).toBeTruthy();
  });

  it(`should have as title 'S.C.I.L.E.R'`, () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app.title).toEqual('S.C.I.L.E.R');
  });

  it('should render title', () => {
    const fixture = TestBed.createComponent(AppComponent);
    fixture.detectChanges();
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector('h1').textContent).toContain('Welcome to S.C.I.L.E.R');
  });
});
