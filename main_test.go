package main

import (
	"testing"
	"time"
)

func TestDateIsToday_Today(t *testing.T) {
	now := time.Now()
	first := now

	result, err := DateIsToday(now, first)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Test if result is true --> it should be
	if !result {
		t.Errorf("expected true for same day, got false")
	}
}

func TestDateIsToday_Yesterday(t *testing.T) {
	now := time.Now()
	first := now.AddDate(0, 0, -1)

	result, err := DateIsToday(now, first)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Should return false --> 
	if result {
		t.Errorf("expected false for yesterday, got true")
	}
}

func TestDateIsToday_Tomorrow(t *testing.T) {
	now := time.Now()
	first := now.AddDate(0, 0, 1)

	result, err := DateIsToday(now, first)
	// Future date should return error
	if err == nil {
		t.Errorf("expected error for future date, got nil")
	}
	if result {
		t.Errorf("expected false for future date, got true")
	}
}
