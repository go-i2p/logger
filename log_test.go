package logger

import (
	"os"
	"testing"
)

// TestBug1NilPointerDereference reproduces and verifies the fix for nil pointer dereference bug
func TestBug1NilPointerDereference(t *testing.T) {
	// Save original state
	originalLog := log
	originalFailFast := failFast

	// Restore state after test
	defer func() {
		log = originalLog
		failFast = originalFailFast
		os.Unsetenv("WARNFAIL_I2P")
		os.Unsetenv("DEBUG_I2P")
	}()

	// Set up conditions that trigger the bug
	log = nil         // Simulate uninitialized logger
	failFast = "true" // Enable fast-fail mode

	// Set environment to enable fast-fail
	os.Setenv("WARNFAIL_I2P", "true")
	os.Setenv("DEBUG_I2P", "debug")

	defer func() {
		if r := recover(); r != nil {
			// This should be a controlled panic, not a nil pointer dereference
			if r == "Logger not initialized but fast-fail mode enabled" {
				t.Logf("Correct controlled panic occurred: %v", r)
			} else {
				t.Errorf("Unexpected panic: %v", r)
			}
		}
	}()

	// This should trigger a controlled panic instead of nil pointer dereference
	warnFatal("test warning")
	t.Fatal("Expected controlled panic, but none occurred")
}

// TestBug1FixVerification verifies that normal operation works after the fix
func TestBug1FixVerification(t *testing.T) {
	// Initialize a new logger properly
	logger := New()

	// Save original state
	originalFailFast := failFast
	defer func() {
		failFast = originalFailFast
		os.Unsetenv("WARNFAIL_I2P")
	}()

	// Test that warnings work normally when failFast is not set
	failFast = ""
	logger.Warn("test warning - should not cause panic")

	// Test should reach here without panicking
	t.Log("Normal warning operation works correctly")
}
