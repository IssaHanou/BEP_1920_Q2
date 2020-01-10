package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGeneralEvent_GetName(t *testing.T) {
	generalEvent := GeneralEvent{
		Name:  "event",
		Rules: nil,
	}

	assert.Equal(t, generalEvent.GetName(), "event")
}

func TestGeneralEvent_GetRules(t *testing.T) {
	rule := new(Rule)

	generalEvent := GeneralEvent{
		Name: "event",
		Rules: []*Rule{
			rule,
		},
	}

	assert.Equal(t, generalEvent.GetRules(), []*Rule{rule})
}

// TODO Test last timer parts: t.Ending() and t.Stop()
func TestTimer_GetTimeLeft(t *testing.T) {
	timer := Timer{
		ID:        "testTimer",
		Duration:  10 * time.Second,
		StartedAt: time.Time{},
		T:         nil,
		State:     "stateIdle",
		Ending:    nil,
		Finish:    false,
	}
	left, state := timer.GetTimeLeft()
	assert.Equal(t, left, 10*time.Second)
	assert.Equal(t, state, "stateIdle")
}

func TestTimer_GetTimeLeft_Active(t *testing.T) {
	timer := Timer{
		ID:        "testTimer",
		Duration:  10 * time.Second,
		StartedAt: time.Time{},
		T:         nil,
		State:     "stateIdle",
		Ending:    nil,
		Finish:    false,
	}
	timer.Start(nil)
	left, state := timer.GetTimeLeft()
	assert.GreaterOrEqual(t, left.Seconds(), (timer.Duration - time.Now().Sub(timer.StartedAt)).Seconds())
	assert.Equal(t, state, "stateActive")
}

func TestTimer_Start(t *testing.T) {
	timer := Timer{
		ID:        "testTimer",
		Duration:  10 * time.Second,
		StartedAt: time.Time{},
		T:         nil,
		State:     "stateIdle",
		Ending:    nil,
		Finish:    false,
	}
	timer.Start(nil)
	assert.Equal(t, timer.State, "stateActive")
}

func TestTimer_Start_False(t *testing.T) {
	timer := Timer{
		ID:        "testTimer",
		Duration:  10 * time.Second,
		StartedAt: time.Time{},
		T:         nil,
		State:     "stateIdle",
		Ending:    nil,
		Finish:    false,
	}
	ok := timer.Start(nil)
	assert.Equal(t, ok, true)
	ok2 := timer.Start(nil)
	assert.Equal(t, ok2, false)
}

func TestTimer_Pause(t *testing.T) {
	timer := Timer{
		ID:        "testTimer",
		Duration:  10 * time.Second,
		StartedAt: time.Time{},
		T:         nil,
		State:     "stateIdle",
		Ending:    nil,
		Finish:    false,
	}
	timer.Start(nil)
	assert.Equal(t, timer.State, "stateActive")
	timer.Pause()
	assert.Equal(t, timer.State, "stateIdle")
}

func TestTimer_Pause_False(t *testing.T) {
	timer := Timer{
		ID:        "testTimer",
		Duration:  10 * time.Second,
		StartedAt: time.Time{},
		T:         nil,
		State:     "stateIdle",
		Ending:    nil,
		Finish:    false,
	}
	ok := timer.Pause()
	assert.Equal(t, timer.State, "stateIdle")
	assert.Equal(t, ok, false)
}

func TestTimer_Stop(t *testing.T) {
	timer := Timer{
		ID:        "testTimer",
		Duration:  10 * time.Second,
		StartedAt: time.Time{},
		T:         nil,
		State:     "stateIdle",
		Ending:    nil,
		Finish:    false,
	}
	timer.Start(nil)
	assert.Equal(t, timer.State, "stateActive")
	timer.Stop()
	assert.Equal(t, timer.State, "stateExpired")
}

func TestTimer_Stop_False(t *testing.T) {
	timer := Timer{
		ID:        "testTimer",
		Duration:  10 * time.Second,
		StartedAt: time.Time{},
		T:         nil,
		State:     "stateIdle",
		Ending:    nil,
		Finish:    false,
	}
	ok := timer.Stop()
	assert.Equal(t, ok, false)
}

func TestPuzzle_GetName(t *testing.T) {
	rule := new(Rule)

	puzzle := Puzzle{
		Event: &GeneralEvent{
			Name: "event",
			Rules: []*Rule{
				rule,
			},
		},
		Hints: nil,
	}

	assert.Equal(t, puzzle.GetName(), "event")
}

func TestPuzzle_GetRules(t *testing.T) {
	rule := new(Rule)
	puzzle := Puzzle{
		Event: &GeneralEvent{
			Name: "event",
			Rules: []*Rule{
				rule,
			},
		},
		Hints: nil,
	}

	assert.Equal(t, puzzle.GetRules(), []*Rule{rule})
}

func Test_CompareWrongComparison(t *testing.T) {
	assert.False(t, compare("a", "a", "unknown"),
		"comparisons should be done on existing options like eq and gte")
}

func Test_CompareWrongComparisonPerType(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckComparisonPerType.json"
	assert.Panics(t,
		func() { ReadFile(filename) },
		"The comparison in a constraint on a condition of type rule may only be a numeric comparator")
}

func Test_compare(t *testing.T) {
	type args struct {
		param1     interface{}
		param2     interface{}
		comparison string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "equal false",
			args: args{
				param1:     float64(1),
				param2:     float64(2),
				comparison: "eq",
			},
			want: false,
		}, {
			name: "equal true",
			args: args{
				param1:     float64(2),
				param2:     float64(2),
				comparison: "eq",
			},
			want: true,
		}, {
			name: "equal true int 2 == float62(2)",
			args: args{
				param1:     2, // int
				param2:     float64(2),
				comparison: "eq",
			},
			want: true,
		}, {
			name: "less then false",
			args: args{
				param1:     float64(1),
				param2:     float64(1),
				comparison: "lt",
			},
			want: false,
		}, {
			name: "less then true",
			args: args{
				param1:     float64(1),
				param2:     float64(2),
				comparison: "lt",
			},
			want: true,
		}, {
			name: "greater then false",
			args: args{
				param1:     float64(1),
				param2:     float64(1),
				comparison: "gt",
			},
			want: false,
		},
		{
			name: "greater then true",
			args: args{
				param1:     float64(2),
				param2:     float64(1),
				comparison: "gt",
			},
			want: true,
		}, {
			name: "less then equal false",
			args: args{
				param1:     float64(2),
				param2:     float64(1),
				comparison: "lte",
			},
			want: false,
		}, {
			name: "less then equal true1",
			args: args{
				param1:     float64(1),
				param2:     float64(1),
				comparison: "lte",
			},
			want: true,
		}, {
			name: "less then equal true2",
			args: args{
				param1:     float64(1),
				param2:     float64(2),
				comparison: "lte",
			},
			want: true,
		}, {
			name: "greater then equal false",
			args: args{
				param1:     float64(1),
				param2:     float64(2),
				comparison: "gte",
			},
			want: false,
		}, {
			name: "greater then equal true1",
			args: args{
				param1:     float64(1),
				param2:     float64(1),
				comparison: "gte",
			},
			want: true,
		}, {
			name: "greater then equal true2",
			args: args{
				param1:     float64(2),
				param2:     float64(1),
				comparison: "gte",
			},
			want: true,
		}, {
			name: "contains false",
			args: args{
				param1:     []float64{1, 2, 3, 5},
				param2:     float64(4),
				comparison: "contains",
			},
			want: false,
		}, {
			name: "contains true",
			args: args{
				param1:     []float64{1, 2, 3, 4, 5},
				param2:     float64(4),
				comparison: "contains",
			},
			want: true,
		}, {
			name: "equals slice true",
			args: args{
				param1:     []float64{1, 2, 3, 4, 5},
				param2:     []float64{1, 2, 3, 4, 5},
				comparison: "eq",
			},
			want: true,
		}, {
			name: "equals slice false",
			args: args{
				param1:     []float64{1, 2, 3, 4, 5},
				param2:     []float64{1, 2, 3, 5},
				comparison: "eq",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compare(tt.args.param1, tt.args.param2, tt.args.comparison); got != tt.want {
				t.Errorf("compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CheckConstraintInputString(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputString.json"
	assert.PanicsWithValue(t,
		"on rule controlSwitch: input type string expected but float64 found as type of value 1",
		func() { ReadFile(filename) },
		"When input is specified as a string, the value of a condition should be a string in order to be able to do a comparison")
}

func Test_CheckConstraintInputStringComparison(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputStringComparison.json"
	assert.PanicsWithValue(t,
		"on rule controlSwitch: comparison lte not allowed on a string",
		func() { ReadFile(filename) },
		"When input is specified as a string, only eq is allowed as comparison")
}

func Test_CheckConstraintInputBool(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputBool.json"
	assert.PanicsWithValue(t,
		"on rule controlSwitch: input type boolean expected but float64 found as type of value 1",
		func() { ReadFile(filename) },
		"When input is specified as a bool, the value of a condition should be a bool in order to be able to do a comparison")
}

func Test_CheckConstraintInputBoolComparison(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputBoolComparison.json"
	assert.PanicsWithValue(t,
		"on rule controlSwitch: comparison lte not allowed on a boolean",
		func() { ReadFile(filename) },
		"When input is specified as a bool, only eq is allowed as comparison")
}

func Test_CheckConstraintInputNumeric(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputNumeric.json"
	assert.PanicsWithValue(t,
		"on rule controlSwitch: input type numeric expected but bool found as type of value true",
		func() { ReadFile(filename) },
		"When input is specified as a numeric, the value of a condition should be a numeric in order to be able to do a comparison")
}

func Test_CheckConstraintInputNumericComparison(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputNumericComparison.json"
	assert.PanicsWithValue(t,
		"on rule controlSwitch: comparison contains not allowed on a numeric",
		func() { ReadFile(filename) },
		"When input is specified as a numeric, only eq, lt, gt, lte, gte are allowed as comparison")
}

func Test_CheckConstraintInputArray(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputArray.json"
	assert.PanicsWithValue(t,
		"on rule controlSwitch: input type array/slice expected but float64 found as type of value 1",
		func() { ReadFile(filename) },
		"When input is specified as a array, the value of a condition should be an array in order to be able to do a comparison")
}

func Test_CheckConstraintInputArrayComparison(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputArrayComparison.json"
	assert.PanicsWithValue(t,
		"on rule controlSwitch: comparison lte not allowed on an array",
		func() { ReadFile(filename) },
		"When input is specified as an array, only contains and eq are allowed as comparison")
}

func Test_CheckConstraintInputCustom(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckConstraintInputCustom.json"
	assert.PanicsWithValue(t,
		"on rule controlSwitch: custom types like: custom, are not yet implemented",
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
	assert.False(t, constraint.Resolve(condition, ReadFile(filename)),
		"Custom condition types are not supported!")
}

func Test_ResolveDeviceFalse(t *testing.T) {
	filename := "../../../resources/testing/test_resolveFalse.json"
	config := ReadFile(filename)
	config.Devices["controlBoard"].Status["mainSwitch"] = false
	config.Devices["controlBoard"].Status["greenSwitch"] = true
	config.Devices["controlBoard"].Status["redSwitch"] = true
	config.Devices["controlBoard"].Status["slider2"] = 30
	assert.False(t, config.Puzzles[1].GetRules()[0].Conditions.Resolve(config))
}

func Test_ResolveDeviceTrue(t *testing.T) {
	filename := "../../../resources/testing/test_resolveTrue.json"
	config := ReadFile(filename)
	config.Devices["controlBoard"].Status["mainSwitch"] = true
	config.Devices["controlBoard"].Status["greenSwitch"] = true
	config.Devices["controlBoard"].Status["redSwitch"] = true
	config.Devices["controlBoard"].Status["slider2"] = 30
	assert.True(t, config.Puzzles[1].GetRules()[0].Conditions.Resolve(config))
}

func Test_ResolveTimerTrue(t *testing.T) {
	filename := "../../../resources/testing/test_resolveTrue.json"
	config := ReadFile(filename)
	config.Timers["timer1"].Finish = true
	assert.True(t, config.GeneralEvents[0].GetRules()[0].Conditions.Resolve(config))
}

func Test_ResolveTimerFalse(t *testing.T) {
	filename := "../../../resources/testing/test_resolveFalse.json"
	config := ReadFile(filename)
	assert.True(t, config.GeneralEvents[0].GetRules()[0].Conditions.Resolve(config))
}

func Test_ResolveRuleTrue(t *testing.T) {
	filename := "../../../resources/testing/test_resolveTrue.json"
	config := ReadFile(filename)
	config.RuleMap["flipSwitch"].Executed = 1
	assert.True(t, config.GeneralEvents[0].GetRules()[0].Conditions.Resolve(config))
}

func Test_ResolveRuleFalse(t *testing.T) {
	filename := "../../../resources/testing/test_resolveFalse.json"
	config := ReadFile(filename)
	assert.False(t, config.GeneralEvents[1].GetRules()[0].Conditions.Resolve(config))
}

func Test_ReadTimer(t *testing.T) {
	filename := "../../../resources/testing/test_config.json"
	config := ReadFile(filename)
	assert.Equal(t, config.Timers["general"].Duration, time.Minute*30)
}

func Test_CheckRuleValue(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckRuleValue.json"
	assert.PanicsWithValue(t,
		"on rule controlSwitch: value type numeric expected but bool found as type of value true",
		func() { ReadFile(filename) },
		"The value in a constraint on a condition of type rule may only be of type numeric")

}

func Test_CheckTimerValue(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckTimerValue.json"
	assert.PanicsWithValue(t,
		"on rule telephoneRings: input type boolean expected but string found as type of value testString",
		func() { ReadFile(filename) },
		"The value in a constraint on a condition of type rule may only be of type bool")

}

func Test_CheckRuleComparison(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckRuleComparison.json"
	assert.PanicsWithValue(t,
		"on rule controlSwitch: comparison contains not allowed on rule",
		func() { ReadFile(filename) },
		"The comparison in a constraint on a condition of type rule may only be a numeric comparator")
}

func Test_CheckTimerComparison(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckTimerComparison.json"
	assert.PanicsWithValue(t,
		"on rule telephoneRings: comparison gte not allowed on a boolean",
		func() { ReadFile(filename) },
		"The comparison in a constraint on a condition of type rule may only be a numeric comparator")
}

func Test_CheckRuleID(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckRuleID.json"
	assert.PanicsWithValue(t,
		"on rule controlSwitch: rule with id non existing not found in map",
		func() { ReadFile(filename) },
		"The rule id is unknown")
}

func Test_CheckTimerID(t *testing.T) {
	filename := "../../../resources/testing/wrong-types/testCheckTimerID.json"
	assert.PanicsWithValue(t,
		"on rule telephoneRings: timer with id timerTest not found in map",
		func() { ReadFile(filename) },
		"The rule id is unknown")
}

func TestNumericToFloatNonNumeric(t *testing.T) {
	assert.Equal(t, float64(0), numericToFloat64("0"),
		"non input or float value should return 0")
}
