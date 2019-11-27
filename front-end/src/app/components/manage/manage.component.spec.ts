import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ManageComponent } from './manage.component';

describe('ManageComponent', () => {
  let component: ManageComponent;
  let fixture: ComponentFixture<ManageComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ManageComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ManageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should render header', () => {
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector('h3').textContent).toContain('Acties');
  });

  it('should have all buttons present', () => {
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.getElementsByClassName('start').item(0).textContent).toContain('Start');
    expect(compiled.getElementsByClassName('test').item(0).textContent).toContain('Test');
    expect(compiled.getElementsByClassName('stop').item(0).textContent).toContain('Stop');
    expect(compiled.getElementsByClassName('reset').item(0).textContent).toContain('Reset');
  });
});
