package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func Test_CheckConstraintInputString(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputString.json"
	assert.PanicsWithValue(t,
		"input type string expected but float64 found as type of value 1",
		func() { ReadFile(filename) },
		"When input is specified as a string, the value of a condition should be a string in order to be able to do a comparison")
}

func Test_CheckConstraintInputStringComparison(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputStringComparison.json"
	assert.PanicsWithValue(t,
		"comparision lte not allowed on a string",
		func() { ReadFile(filename) },
		"When input is specified as a string, only eq is allowed as comparison")
}

func Test_CheckConstraintInputBool(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputBool.json"
	assert.PanicsWithValue(t,
		"input type boolean expected but float64 found as type of value 1",
		func() { ReadFile(filename) },
		"When input is specified as a bool, the value of a condition should be a bool in order to be able to do a comparison")
}

func Test_CheckConstraintInputBoolComparison(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputBoolComparison.json"
	assert.PanicsWithValue(t,
		"comparision lte not allowed on a boolean",
		func() { ReadFile(filename) },
		"When input is specified as a bool, only eq is allowed as comparison")
}

func Test_CheckConstraintInputNumeric(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputNumeric.json"
	assert.PanicsWithValue(t,
		"input type numeric expected but bool found as type of value true",
		func() { ReadFile(filename) },
		"When input is specified as a numeric, the value of a condition should be a numeric in order to be able to do a comparison")
}

func Test_CheckConstraintInputNumericComparison(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputNumericComparison.json"
	assert.PanicsWithValue(t,
		"comparision contains not allowed on a numeric",
		func() { ReadFile(filename) },
		"When input is specified as a numeric, only eq, lt, gt, lte, gte are allowed as comparison")
}

func Test_CheckConstraintInputArray(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputArray.json"
	assert.PanicsWithValue(t,
		"input type array/slice expected but float64 found as type of value 1",
		func() { ReadFile(filename) },
		"When input is specified as a array, the value of a condition should be an array in order to be able to do a comparison")
}

func Test_CheckConstraintInputArrayComparison(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputArrayComparison.json"
	assert.PanicsWithValue(t,
		"comparision lte not allowed on an array",
		func() { ReadFile(filename) },
		"When input is specified as an array, only contains and eq are allowed as comparison")
}

func Test_CheckConstraintInputCustom(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputCustom.json"
	assert.PanicsWithValue(t,
		"custom types like: custom, are not yet implemented",
		func() { ReadFile(filename) },
		"Custom types are not supported yet")
}

func Test_CheckResolve(t *testing.T) {
	constraint := Constraint{
		Comparison:  "eq",
		ComponentID: "redSwitch",
		Value:       true,
	}

	condition := Condition{
		Type:        "custom",
		TypeID:      "controlBoard",
		Constraints: constraint,
	}
	filename := "../../../resources/testing/test_config.json"
	assert.PanicsWithValue(t,
		fmt.Sprintf("cannot resolve constraint %v because condition.type is an unknown type", constraint),
		func() { constraint.Resolve(condition, ReadFile(filename)) },
		"Custom condition types are not supported!")
}
