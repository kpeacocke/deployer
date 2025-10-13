package main

import (
	"path/filepath"
	"testing"
)

func TestLoadState(t *testing.T) {
	// Test loading non-existent state file (should return defaults)
	tempDir := t.TempDir()
	statePath := filepath.Join(tempDir, "nonexistent.yaml")

	state, err := LoadState(statePath)
	if err != nil {
		t.Fatalf("Failed to load non-existent state: %v", err)
	}

	if state.ActiveSlot != "blue" {
		t.Errorf("Expected default active slot 'blue', got '%s'", state.ActiveSlot)
	}

	if state.BlueVersion != "" {
		t.Errorf("Expected empty blue version, got '%s'", state.BlueVersion)
	}
}

func TestSaveAndLoadState(t *testing.T) {
	tempDir := t.TempDir()
	statePath := filepath.Join(tempDir, "state.yaml")

	// Create state
	originalState := &DeploymentState{
		ActiveSlot:   "green",
		BlueVersion:  "v1.0.0",
		GreenVersion: "v1.1.0",
	}

	// Save state
	if err := originalState.SaveState(statePath); err != nil {
		t.Fatalf("Failed to save state: %v", err)
	}

	// Load state
	loadedState, err := LoadState(statePath)
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	// Compare
	if loadedState.ActiveSlot != originalState.ActiveSlot {
		t.Errorf("Expected active slot '%s', got '%s'", originalState.ActiveSlot, loadedState.ActiveSlot)
	}

	if loadedState.BlueVersion != originalState.BlueVersion {
		t.Errorf("Expected blue version '%s', got '%s'", originalState.BlueVersion, loadedState.BlueVersion)
	}

	if loadedState.GreenVersion != originalState.GreenVersion {
		t.Errorf("Expected green version '%s', got '%s'", originalState.GreenVersion, loadedState.GreenVersion)
	}
}

func TestGetInactiveSlot(t *testing.T) {
	tests := []struct {
		activeSlot       string
		expectedInactive string
	}{
		{"blue", "green"},
		{"green", "blue"},
	}

	for _, test := range tests {
		state := &DeploymentState{ActiveSlot: test.activeSlot}
		inactive := state.GetInactiveSlot()

		if inactive != test.expectedInactive {
			t.Errorf("For active slot '%s', expected inactive '%s', got '%s'",
				test.activeSlot, test.expectedInactive, inactive)
		}
	}
}

func TestSwitchSlot(t *testing.T) {
	state := &DeploymentState{ActiveSlot: "blue"}

	state.SwitchSlot()
	if state.ActiveSlot != "green" {
		t.Errorf("Expected active slot 'green' after switch, got '%s'", state.ActiveSlot)
	}

	state.SwitchSlot()
	if state.ActiveSlot != "blue" {
		t.Errorf("Expected active slot 'blue' after second switch, got '%s'", state.ActiveSlot)
	}
}
