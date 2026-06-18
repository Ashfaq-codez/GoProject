package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	now := time.Now()
	
	// Test someone born exactly 20 years ago today
	dob20 := now.AddDate(-20, 0, 0)
	if age := CalculateAge(dob20); age != 20 {
		t.Errorf("Expected 20, got %d", age)
	}

	// Test someone born 20 years ago, but their birthday is tomorrow
	dobAlmost20 := now.AddDate(-20, 0, 1)
	if age := CalculateAge(dobAlmost20); age != 19 {
		t.Errorf("Expected 19, got %d", age)
	}
}