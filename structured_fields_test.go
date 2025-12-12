package logger

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestStructuredFieldsOutput(t *testing.T) {
	// Create a new logger for testing
	log := New()

	// Create a buffer to capture output
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	log.SetLevel(InfoLevel)

	// Configure formatter
	fmtter := &TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		ForceColors:            false,
		DisableColors:          true,
		DisableQuote:           true,
		DisableTimestamp:       false,
		DisableSorting:         false,
		DisableLevelTruncation: false,
		QuoteEmptyFields:       true,
	}
	log.SetFormatter(fmtter)

	// Test 1: Simple structured logging
	t.Run("SimpleFields", func(t *testing.T) {
		buf.Reset()
		log.WithFields(Fields{
			"test":   "value",
			"number": 42,
		}).Info("test message")

		output := buf.String()
		if !strings.Contains(output, "test=value") {
			t.Errorf("Expected output to contain 'test=value', got: %s", output)
		}
		if !strings.Contains(output, "number=42") {
			t.Errorf("Expected output to contain 'number=42', got: %s", output)
		}
		if !strings.Contains(output, "msg=\"test message\"") && !strings.Contains(output, "msg=test message") {
			t.Errorf("Expected output to contain message, got: %s", output)
		}
	})

	// Test 2: I2CP payload size example from the bug report
	t.Run("I2CPPayloadFields", func(t *testing.T) {
		buf.Reset()
		payloadLen := 500000
		MaxPayloadSize := 1048576
		log.WithFields(Fields{
			"at":         "i2cp.ReadMessage",
			"payloadLen": payloadLen,
			"maxAllowed": MaxPayloadSize,
			"exceeded":   payloadLen - MaxPayloadSize,
		}).Error("payload_size_exceeded_max")

		output := buf.String()
		if !strings.Contains(output, "at=i2cp.ReadMessage") && !strings.Contains(output, "at=\"i2cp.ReadMessage\"") {
			t.Errorf("Expected output to contain 'at=i2cp.ReadMessage', got: %s", output)
		}
		if !strings.Contains(output, "payloadLen=500000") {
			t.Errorf("Expected output to contain 'payloadLen=500000', got: %s", output)
		}
		if !strings.Contains(output, "maxAllowed=1048576") {
			t.Errorf("Expected output to contain 'maxAllowed=1048576', got: %s", output)
		}
		if !strings.Contains(output, "exceeded=-548576") {
			t.Errorf("Expected output to contain 'exceeded=-548576', got: %s", output)
		}
	})

	// Test 3: WithField chaining
	t.Run("ChainedFields", func(t *testing.T) {
		buf.Reset()
		log.WithField("first", "value1").
			WithField("second", "value2").
			Info("chained fields")

		output := buf.String()
		if !strings.Contains(output, "first=value1") {
			t.Errorf("Expected output to contain 'first=value1', got: %s", output)
		}
		if !strings.Contains(output, "second=value2") {
			t.Errorf("Expected output to contain 'second=value2', got: %s", output)
		}
	})

	// Test 4: Warn and Error with fields preserve failFast behavior
	t.Run("WarnWithFields", func(t *testing.T) {
		buf.Reset()
		// Should not panic since WARNFAIL_I2P is not set
		log.WithFields(Fields{
			"key": "value",
		}).Warn("warning message")

		output := buf.String()
		if !strings.Contains(output, "key=value") {
			t.Errorf("Expected output to contain 'key=value', got: %s", output)
		}
	})
}

func TestGlobalLoggerWithFields(t *testing.T) {
	// Test the global package-level functions
	os.Setenv("DEBUG_I2P", "info")
	InitializeGoI2PLogger()

	log := GetGoI2PLogger()
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	log.SetLevel(InfoLevel)

	// Test WithFields global function
	WithFields(Fields{
		"global": "test",
		"number": 99,
	}).Info("global logger test")

	output := buf.String()
	if !strings.Contains(output, "global=test") {
		t.Errorf("Expected global logger output to contain 'global=test', got: %s", output)
	}
	if !strings.Contains(output, "number=99") {
		t.Errorf("Expected global logger output to contain 'number=99', got: %s", output)
	}

	// Cleanup
	os.Unsetenv("DEBUG_I2P")
}
