package utils

import (
	"testing"
	"time"
)

func mustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}

func TestParseRaceResult(t *testing.T) {
	tests := []struct {
		input    string
		wantLaps int
		wantDur  time.Duration
		wantErr  bool
	}{
		{"11/34.234", 11, mustParseDuration("34.234s"), false},
		{"11/1'34.234", 11, mustParseDuration("1m34.234s"), false},
		{"110/20'69.1234", 110, mustParseDuration("20m69.1234s"), false},
		{"-1/30.5", -1, mustParseDuration("30.5s"), false},
		{"-1/2'56.521", -1, mustParseDuration("2m56.521s"), false},
		{"invalid", 0, 0, true},
	}

	for _, tt := range tests {
		laps, dur, err := ParseRaceResult(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseRaceResult(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if laps != tt.wantLaps || dur != tt.wantDur {
			t.Errorf("ParseRaceResult(%q) = %d, %v, want %d, %v", tt.input, laps, dur, tt.wantLaps, tt.wantDur)
		}
	}
}

func TestParseGap(t *testing.T) {
	tests := []struct {
		input   string
		wantDur time.Duration
		wantErr bool
	}{
		{"+3.45", mustParseDuration("3.45s"), false},
		{"+0.45", mustParseDuration("0.45s"), false},
		{"+110.987", mustParseDuration("110.987s"), false},
		{"+0.001", mustParseDuration("0.001s"), false},
		{"3.45", 0, true}, // Missing '+' prefix
		{"+invalid", 0, true},
	}

	for _, tt := range tests {
		dur, err := ParseGap(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseGap(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if dur != tt.wantDur {
			t.Errorf("ParseGap(%q) = %v, want %v", tt.input, dur, tt.wantDur)
		}
	}
}

func TestParseLap(t *testing.T) {
	tests := []struct {
		input   string
		wantDur time.Duration
		wantLap int
		wantErr bool
	}{
		{"9.000", mustParseDuration("9s"), 0, false},
		{"9.001", mustParseDuration("9.001s"), 0, false},
		{"11.345", mustParseDuration("11.345s"), 0, false},
		{"23.543[9]", mustParseDuration("23.543s"), 9, false},
		{"10.0[100]", mustParseDuration("10s"), 100, false},
		{"bad", 0, 0, true},
	}

	for _, tt := range tests {
		dur, lap, err := ParseLap(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseLap(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if dur != tt.wantDur || lap != tt.wantLap {
			t.Errorf("ParseLap(%q) = %v, %d, want %v, %d", tt.input, dur, lap, tt.wantDur, tt.wantLap)
		}
	}
}
