package main

import (
	"bytes"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

type renderNotifyingWriter struct {
	mu       sync.Mutex
	buf      bytes.Buffer
	rendered chan struct{}
	once     sync.Once
}

func (w *renderNotifyingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	n, err := w.buf.Write(p)
	if strings.Contains(w.buf.String(), "\x1b[0m") {
		w.once.Do(func() { close(w.rendered) })
	}
	return n, err
}

func (w *renderNotifyingWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.buf.String()
}

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

func TestRunHelp(t *testing.T) {
	for _, args := range [][]string{{"help"}, {"-h"}, {"--help"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var out bytes.Buffer
			if err := run(args, &out); err != nil {
				t.Fatalf("run(%v) error = %v", args, err)
			}
			if got := out.String(); got != usage {
				t.Fatalf("run(%v) output = %q, want %q", args, got, usage)
			}
		})
	}
}

func TestCLIHelpExitsSuccessfully(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	if got := cli([]string{"help"}, &out, &errOut); got != 0 {
		t.Fatalf("cli() = %d, want 0", got)
	}
	if got := out.String(); got != usage {
		t.Fatalf("cli() stdout = %q, want %q", got, usage)
	}
	if got := errOut.String(); got != "" {
		t.Fatalf("cli() stderr = %q, want empty stderr", got)
	}
}

func TestCLIInvalidInputExitsWithUsage(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	if got := cli([]string{"nope"}, &out, &errOut); got != 1 {
		t.Fatalf("cli() = %d, want 1", got)
	}
	if got := out.String(); got != "" {
		t.Fatalf("cli() stdout = %q, want empty stdout", got)
	}
	got := errOut.String()
	for _, want := range []string{"invalid duration", usage} {
		if !strings.Contains(got, want) {
			t.Fatalf("cli() stderr = %q, want %q", got, want)
		}
	}
}

func TestRunRejectsInvalidInput(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"not-a-duration"}, &out); err == nil {
		t.Fatal("run() expected error")
	}
	if got := out.String(); got != "" {
		t.Fatalf("run() output = %q, want empty output", got)
	}
}

func TestRunCompletesShortTimer(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"1ns"}, &out); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{renderLarge("00:00"), "Time's up!", "\x1b[?25h"} {
		if !strings.Contains(got, want) {
			t.Fatalf("run() output = %q, want %q", got, want)
		}
	}
}

func TestWantsHelp(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{name: "help", args: []string{"help"}, want: true},
		{name: "short flag", args: []string{"-h"}, want: true},
		{name: "long flag", args: []string{"--help"}, want: true},
		{name: "empty", args: nil, want: false},
		{name: "duration", args: []string{"10s"}, want: false},
		{name: "extra arg", args: []string{"help", "10s"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wantsHelp(tt.args); got != tt.want {
				t.Fatalf("wantsHelp() = %v, want %v", got, tt.want)
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

func TestRenderLargeIgnoresUnsupportedCharacters(t *testing.T) {
	if got, want := renderLarge("x"), "\n\n\n\n\n"; got != want {
		t.Fatalf("renderLarge() = %q, want %q", got, want)
	}
}

func TestRenderCountdown(t *testing.T) {
	var out bytes.Buffer
	renderCountdown(10*time.Second, &out)

	got := out.String()
	for _, want := range []string{"\x1b[2J\x1b[H", "\x1b[38;5;208m", renderLarge("00:10"), "\x1b[0m"} {
		if !strings.Contains(got, want) {
			t.Fatalf("renderCountdown() output = %q, want %q", got, want)
		}
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

func TestCountdownCompletionRestoresCursor(t *testing.T) {
	var out bytes.Buffer
	countdown(1*time.Nanosecond, &out, make(chan os.Signal))

	got := out.String()
	for _, want := range []string{renderLarge("00:00"), "Time's up!", "\x1b[?25h"} {
		if !strings.Contains(got, want) {
			t.Fatalf("countdown() output = %q, want %q", got, want)
		}
	}
}

func TestCountdownSubSecondTimerCompletes(t *testing.T) {
	var out bytes.Buffer
	countdown(time.Millisecond, &out, make(chan os.Signal))

	got := out.String()
	for _, want := range []string{renderLarge("00:00"), "Time's up!", "\x1b[?25h"} {
		if !strings.Contains(got, want) {
			t.Fatalf("countdown() output = %q, want %q", got, want)
		}
	}
}

func TestCountdownCancellationDuringSleepRestoresCursor(t *testing.T) {
	interrupts := make(chan os.Signal, 1)
	out := &renderNotifyingWriter{rendered: make(chan struct{})}
	done := make(chan struct{})

	go func() {
		countdown(time.Hour, out, interrupts)
		close(done)
	}()

	select {
	case <-out.rendered:
	case <-time.After(time.Second):
		t.Fatal("countdown() did not render before timeout")
	}

	interrupts <- os.Interrupt

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("countdown() did not cancel before timeout")
	}

	got := out.String()
	for _, want := range []string{"Timer cancelled.", "\x1b[?25h"} {
		if !strings.Contains(got, want) {
			t.Fatalf("countdown() output = %q, want %q", got, want)
		}
	}
}
