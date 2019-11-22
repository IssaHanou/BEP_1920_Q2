package config

import (
	"../../config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFromJson(t *testing.T) {
	json := []byte(`{
			"name": "My Awesome Escape",
			"duration": "00:30:00"
		}`)
	result := config.GetFromJson(json)
	expected := "Escape room My Awesome Escape should be solved within 30 minutes"
	assert.Equal(t, result, expected, "JSON should be properly converted to string")
}
