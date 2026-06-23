package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRayColorOutput(t *testing.T) {
	c := NewVec3(0.5, 0.7, 1.0)
	var buf bytes.Buffer
	c.WriteColor(&buf)

	output := buf.String()
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("Output should end with newline")
	}

	// Parse the bytes
	var r, g, b int
	_, err := fmt.Sscanf(output, "%d %d %d", &r, &g, &b)
	if err != nil {
		t.Errorf("Failed to parse color output: %v", err)
	}

	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		t.Errorf("Color values out of range: %d %d %d", r, g, b)
	}
}

func TestColorToByte(t *testing.T) {
	tests := []struct {
		input    float64
		expected int
	}{
		{0.0, 0},
		{1.0, 255},
		{0.5, 181},
	}

	for _, tt := range tests {
		result := colorToByte(tt.input)
		if result != tt.expected {
			t.Errorf("colorToByte(%f) = %d, want %d", tt.input, result, tt.expected)
		}
	}
}
