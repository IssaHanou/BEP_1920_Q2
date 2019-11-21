package config

import (
	"../../src/config"
	"testing"
)

func TestGetFromJson(t *testing.T) {
	json := []byte(`{
			"name": "My Awesome Escape",
			"duration": "00:30:00"
		}`)
	result := config.GetFromJson(json)
	expected := "Escape room My Awesome Escape should be solved within 30 minutes"
	if result != expected {
		t.Error("Excpected: "+expected+"; but was:", result)
	}
}
