package main

import "testing"

func TestParseValue(t *testing.T) {
	if value := parseValue("K0192224409", 6, 3, 4); value != 29.2 {
		t.Errorf("Expected %v, received %v", 29.2, value)
	}
}

func TestParseValueOneDigit(t *testing.T) {
	if value := parseValue("K01934088F2", 6, 3, 4); value != 9.3 {
		t.Errorf("Expected %v, received %v", 9.3, value)
	}
}
