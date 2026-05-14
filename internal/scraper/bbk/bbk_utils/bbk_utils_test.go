package bbk_utils

import (
	"testing"
	"time"
)

func ptr[T any](v T) *T {
	return &v
}

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
		wantLaps *int
		wantDur  *time.Duration
		wantErr  bool
	}{
		{"11/34.234", ptr(11), ptr(mustParseDuration("34.234s")), false},
		{"11/1'34.234", ptr(11), ptr(mustParseDuration("1m34.234s")), false},
		{"110/20'69.1234", ptr(110), ptr(mustParseDuration("20m69.1234s")), false},
		{"-1/30.5", ptr(-1), ptr(mustParseDuration("30.5s")), false},
		{"-1/2'56.521", ptr(-1), ptr(mustParseDuration("2m56.521s")), false},
		{"invalid", nil, nil, true},
	}

	for _, tt := range tests {
		laps, dur, err := ParseRaceResult(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("incorrect error running ParseRaceResult(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}

		// Compare pointers: both nil or both pointing to equal values
		if (laps == nil) != (tt.wantLaps == nil) || (laps != nil && *laps != *tt.wantLaps) {
			t.Errorf("incorrect laps returned by ParseRaceResult(%q) laps = %v, want %v", tt.input, laps, tt.wantLaps)
		}
		if (dur == nil) != (tt.wantDur == nil) || (dur != nil && *dur != *tt.wantDur) {
			t.Errorf("incorrect duration returned by ParseRaceResult(%q) dur = %v, want %v", tt.input, dur, tt.wantDur)
		}
	}
}

func TestParseGap(t *testing.T) {
	tests := []struct {
		input   string
		wantDur *time.Duration
		wantErr bool
	}{
		{"+3.45", ptr(mustParseDuration("3.45s")), false},
		{"+0.45", ptr(mustParseDuration("0.45s")), false},
		{"+110.987", ptr(mustParseDuration("110.987s")), false},
		{"+0.001", ptr(mustParseDuration("0.001s")), false},
		{"3.45", nil, true},
		{"+invalid", nil, true},
	}

	for _, tt := range tests {
		dur, err := ParseGap(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("incorrect error returned by ParseGap(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if (dur == nil) != (tt.wantDur == nil) || (dur != nil && *dur != *tt.wantDur) {
			t.Errorf("incorrect duration returned by ParseGap(%q) = %v, want %v", tt.input, dur, tt.wantDur)
		}
	}
}

func TestParseLap(t *testing.T) {
	tests := []struct {
		input   string
		wantDur *time.Duration
		wantLap *int
		wantErr bool
	}{
		{"9.000", ptr(mustParseDuration("9s")), nil, false},
		{"11.345", ptr(mustParseDuration("11.345s")), nil, false},
		{"23.543[9]", ptr(mustParseDuration("23.543s")), ptr(9), false},
		{"10.0[100]", ptr(mustParseDuration("10s")), ptr(100), false},
		{"bad", nil, nil, true},
	}

	for _, tt := range tests {
		dur, lap, err := ParseLap(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("incorrect errror returned by ParseLap(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if (dur == nil) != (tt.wantDur == nil) || (dur != nil && *dur != *tt.wantDur) {
			t.Errorf("incorrect duration returned by ParseLap(%q) dur = %v, want %v", tt.input, dur, tt.wantDur)
		}
		if (lap == nil) != (tt.wantLap == nil) || (lap != nil && *lap != *tt.wantLap) {
			t.Errorf("incorrect lap returned by ParseLap(%q) lap = %v, want %v", tt.input, lap, tt.wantLap)
		}
	}
}
