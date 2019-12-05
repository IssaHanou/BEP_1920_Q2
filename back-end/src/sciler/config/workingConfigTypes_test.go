package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestAndCondition_CheckConstraints(t *testing.T) {
//	type fields struct {
//		logics []LogicalCondition
//	}
//	type args struct {
//		config WorkingConfig
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			and := AndCondition{
//				logics: tt.fields.logics,
//			}
//			if err := and.CheckConstraints(tt.args.config); (err != nil) != tt.wantErr {
//				t.Errorf("CheckConstraints() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestAndCondition_Resolve(t *testing.T) {
//	type fields struct {
//		logics []LogicalCondition
//	}
//	type args struct {
//		config WorkingConfig
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			and := AndCondition{
//				logics: tt.fields.logics,
//			}
//			if got := and.Resolve(tt.args.config); got != tt.want {
//				t.Errorf("Resolve() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestAndConstraint_CheckConstraints(t *testing.T) {
//	type fields struct {
//		logics []LogicalConstraint
//	}
//	type args struct {
//		condition Condition
//		config    WorkingConfig
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			and := AndConstraint{
//				logics: tt.fields.logics,
//			}
//			if err := and.CheckConstraints(tt.args.condition, tt.args.config); (err != nil) != tt.wantErr {
//				t.Errorf("CheckConstraints() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestAndConstraint_Resolve(t *testing.T) {
//	type fields struct {
//		logics []LogicalConstraint
//	}
//	type args struct {
//		condition Condition
//		config    WorkingConfig
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			and := AndConstraint{
//				logics: tt.fields.logics,
//			}
//			if got := and.Resolve(tt.args.condition, tt.args.config); got != tt.want {
//				t.Errorf("Resolve() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestCondition_CheckConstraints(t *testing.T) {
//	type fields struct {
//		Type        string
//		TypeID      string
//		Constraints LogicalConstraint
//	}
//	type args struct {
//		config WorkingConfig
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			condition := Condition{
//				Type:        tt.fields.Type,
//				TypeID:      tt.fields.TypeID,
//				Constraints: tt.fields.Constraints,
//			}
//			if err := condition.CheckConstraints(tt.args.config); (err != nil) != tt.wantErr {
//				t.Errorf("CheckConstraints() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestCondition_Resolve(t *testing.T) {
//	type fields struct {
//		Type        string
//		TypeID      string
//		Constraints LogicalConstraint
//	}
//	type args struct {
//		config WorkingConfig
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			condition := Condition{
//				Type:        tt.fields.Type,
//				TypeID:      tt.fields.TypeID,
//				Constraints: tt.fields.Constraints,
//			}
//			if got := condition.Resolve(tt.args.config); got != tt.want {
//				t.Errorf("Resolve() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestConstraint_CheckConstraints(t *testing.T) {
//	type fields struct {
//		Comparison  string
//		ComponentID string
//		Value       interface{}
//	}
//	type args struct {
//		condition Condition
//		config    WorkingConfig
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			constraint := Constraint{
//				Comparison:  tt.fields.Comparison,
//				ComponentID: tt.fields.ComponentID,
//				Value:       tt.fields.Value,
//			}
//			if err := constraint.CheckConstraints(tt.args.condition, tt.args.config); (err != nil) != tt.wantErr {
//				t.Errorf("CheckConstraints() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestConstraint_Resolve(t *testing.T) {
//	type fields struct {
//		Comparison  string
//		ComponentID string
//		Value       interface{}
//	}
//	type args struct {
//		condition Condition
//		config    WorkingConfig
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			constraint := Constraint{
//				Comparison:  tt.fields.Comparison,
//				ComponentID: tt.fields.ComponentID,
//				Value:       tt.fields.Value,
//			}
//			if got := constraint.Resolve(tt.args.condition, tt.args.config); got != tt.want {
//				t.Errorf("Resolve() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestGeneralEvent_GetName(t *testing.T) {
	generalEvent := GeneralEvent{
		Name:  "event",
		Rules: nil,
	}

	assert.Equal(t, generalEvent.GetName(), "event")
}

func TestGeneralEvent_GetRules(t *testing.T) {
	generalEvent := GeneralEvent{
		Name:  "event",
		Rules: make([]Rule, 2),
	}

	assert.Equal(t, generalEvent.GetRules(), make([]Rule, 2))
}

//func TestOrCondition_CheckConstraints(t *testing.T) {
//	type fields struct {
//		logics []LogicalCondition
//	}
//	type args struct {
//		config WorkingConfig
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			or := OrCondition{
//				logics: tt.fields.logics,
//			}
//			if err := or.CheckConstraints(tt.args.config); (err != nil) != tt.wantErr {
//				t.Errorf("CheckConstraints() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestOrCondition_Resolve(t *testing.T) {
//	type fields struct {
//		logics []LogicalCondition
//	}
//	type args struct {
//		config WorkingConfig
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			or := OrCondition{
//				logics: tt.fields.logics,
//			}
//			if got := or.Resolve(tt.args.config); got != tt.want {
//				t.Errorf("Resolve() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestOrConstraint_CheckConstraints(t *testing.T) {
//	type fields struct {
//		logics []LogicalConstraint
//	}
//	type args struct {
//		condition Condition
//		config    WorkingConfig
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			or := OrConstraint{
//				logics: tt.fields.logics,
//			}
//			if err := or.CheckConstraints(tt.args.condition, tt.args.config); (err != nil) != tt.wantErr {
//				t.Errorf("CheckConstraints() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestOrConstraint_Resolve(t *testing.T) {
//	type fields struct {
//		logics []LogicalConstraint
//	}
//	type args struct {
//		condition Condition
//		config    WorkingConfig
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			or := OrConstraint{
//				logics: tt.fields.logics,
//			}
//			if got := or.Resolve(tt.args.condition, tt.args.config); got != tt.want {
//				t.Errorf("Resolve() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestPuzzle_GetName(t *testing.T) {
	puzzle := Puzzle{
		Event: GeneralEvent{
			Name:  "event",
			Rules: make([]Rule, 2)},
		Hints: nil,
	}

	assert.Equal(t, puzzle.GetName(), "event")
}

func TestPuzzle_GetRules(t *testing.T) {
	puzzle := Puzzle{
		Event: GeneralEvent{
			Name:  "event",
			Rules: make([]Rule, 2)},
		Hints: nil,
	}

	assert.Equal(t, puzzle.GetRules(), make([]Rule, 2))
}

func Test_CompareWrongComparison(t *testing.T) {
	assert.PanicsWithValue(t,
		"Cannot compare on: unknown",
		func() { compare("a", "a", "unknown") },
		"comparisons should be done on existing options like eq and gte")
}

func Test_compare(t *testing.T) {
	type args struct {
		param1      interface{}
		param2      interface{}
		comparision string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "equal false",
			args: args{
				param1:      float64(1),
				param2:      float64(2),
				comparision: "eq",
			},
			want: false,
		}, {
			name: "equal true",
			args: args{
				param1:      float64(2),
				param2:      float64(2),
				comparision: "eq",
			},
			want: true,
		}, {
			name: "less then false",
			args: args{
				param1:      float64(1),
				param2:      float64(1),
				comparision: "lt",
			},
			want: false,
		}, {
			name: "less then true",
			args: args{
				param1:      float64(1),
				param2:      float64(2),
				comparision: "lt",
			},
			want: true,
		}, {
			name: "greater then false",
			args: args{
				param1:      float64(1),
				param2:      float64(1),
				comparision: "gt",
			},
			want: false,
		},
		{
			name: "greater then true",
			args: args{
				param1:      float64(2),
				param2:      float64(1),
				comparision: "gt",
			},
			want: true,
		}, {
			name: "less then equal false",
			args: args{
				param1:      float64(2),
				param2:      float64(1),
				comparision: "lte",
			},
			want: false,
		}, {
			name: "less then equal true1",
			args: args{
				param1:      float64(1),
				param2:      float64(1),
				comparision: "lte",
			},
			want: true,
		}, {
			name: "less then equal true2",
			args: args{
				param1:      float64(1),
				param2:      float64(2),
				comparision: "lte",
			},
			want: true,
		}, {
			name: "greater then equal false",
			args: args{
				param1:      float64(1),
				param2:      float64(2),
				comparision: "gte",
			},
			want: false,
		}, {
			name: "greater then equal true1",
			args: args{
				param1:      float64(1),
				param2:      float64(1),
				comparision: "gte",
			},
			want: true,
		}, {
			name: "greater then equal true2",
			args: args{
				param1:      float64(2),
				param2:      float64(1),
				comparision: "gte",
			},
			want: true,
		}, {
			name: "contains false",
			args: args{
				param1:      []float64{1, 2, 3, 5},
				param2:      float64(4),
				comparision: "contains",
			},
			want: false,
		}, {
			name: "contains true",
			args: args{
				param1:      []float64{1, 2, 3, 4, 5},
				param2:      float64(4),
				comparision: "contains",
			},
			want: true,
		}, {
			name: "equals slice true",
			args: args{
				param1:      []float64{1, 2, 3, 4, 5},
				param2:      []float64{1, 2, 3, 4, 5},
				comparision: "eq",
			},
			want: true,
		}, {
			name: "equals slice false",
			args: args{
				param1:      []float64{1, 2, 3, 4, 5},
				param2:      []float64{1, 2, 3, 5},
				comparision: "eq",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compare(tt.args.param1, tt.args.param2, tt.args.comparision); got != tt.want {
				t.Errorf("compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	type args struct {
		list    []interface{}
		element interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.list, tt.args.element); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
