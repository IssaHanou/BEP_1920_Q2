package config

import (
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
