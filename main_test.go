package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want time.Duration
	}{
		{name: "seconds", args: []string{"10s"}, want: 10 * time.Second},
		{name: "minutes", args: []string{"6m"}, want: 6 * time.Minute},
		{name: "combined", args: []string{"6m30s"}, want: 6*time.Minute + 30*time.Second},
		{name: "hours", args: []string{"1h"}, want: time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs(tt.args)
			if err != nil {
				t.Fatalf("parseArgs() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("parseArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseArgsRejectsInvalidInput(t *testing.T) {
	tests := [][]string{
		{},
		{"10"},
		{"0s"},
		{"-1s"},
		{"10s", "20s"},
	}

	for _, args := range tests {
		if _, err := parseArgs(args); err == nil {
			t.Fatalf("parseArgs(%v) expected error", args)
		}
	}
}

func TestFormatCompact(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
		want string
	}{
		{name: "seconds", d: 10 * time.Second, want: "00:10"},
		{name: "minutes", d: 6 * time.Minute, want: "06:00"},
		{name: "combined", d: 6*time.Minute + 30*time.Second, want: "06:30"},
		{name: "hour", d: time.Hour, want: "1:00:00"},
		{name: "multi hour", d: 12*time.Hour + 5*time.Minute, want: "12:05:00"},
		{name: "negative", d: -1 * time.Second, want: "00:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatCompact(tt.d); got != tt.want {
				t.Fatalf("formatCompact() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRenderLarge(t *testing.T) {
	got := renderLarge("01:23")
	lines := strings.Split(strings.TrimSuffix(got, "\n"), "\n")

	if len(lines) != 5 {
		t.Fatalf("renderLarge() produced %d lines, want 5", len(lines))
	}
	if !strings.Contains(got, "███") {
		t.Fatalf("renderLarge() = %q, want seven-segment blocks", got)
	}
	if !strings.Contains(got, "  █  ") {
		t.Fatalf("renderLarge() = %q, want colon or one segment", got)
	}
}

func TestCountdownCancellationRestoresCursor(t *testing.T) {
	interrupts := make(chan os.Signal)
	close(interrupts)

	var out bytes.Buffer
	countdown(time.Hour, &out, interrupts)

	got := out.String()
	for _, want := range []string{"\x1b[?25l", "Timer cancelled.", "\x1b[?25h"} {
		if !strings.Contains(got, want) {
			t.Fatalf("countdown() output = %q, want %q", got, want)
		}
	}
}
