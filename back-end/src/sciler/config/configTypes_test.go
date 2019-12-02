package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCondition_GetID(t *testing.T) {
	constraint := ConstraintInfo{"comp": "eq", "value": 4, "component_id": "turn"}
	cond := Condition{"device", "telephone", []ConstraintInfo{constraint}, "ruleID"}
	cond2 := Condition{"device", "telephone", []ConstraintInfo{constraint}, "ruleID"}
	assert.Equal(t, cond.GetID(), cond2.GetID(), "IDs of conditions with same values should be equal")
}
