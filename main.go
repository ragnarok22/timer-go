package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const usage = `Usage: timer <duration>

Examples:
  timer 10s
  timer 6m
  timer 6m30s
  timer 1h
`

var digitPatterns = map[rune][5]string{
	'0': {" ███ ", "█   █", "█   █", "█   █", " ███ "},
	'1': {"  █  ", " ██  ", "  █  ", "  █  ", " ███ "},
	'2': {" ███ ", "█   █", "   █ ", "  █  ", "█████"},
	'3': {"████ ", "    █", " ███ ", "    █", "████ "},
	'4': {"█   █", "█   █", "█████", "    █", "    █"},
	'5': {"█████", "█    ", "████ ", "    █", "████ "},
	'6': {" ███ ", "█    ", "████ ", "█   █", " ███ "},
	'7': {"█████", "    █", "   █ ", "  █  ", "  █  "},
	'8': {" ███ ", "█   █", " ███ ", "█   █", " ███ "},
	'9': {" ███ ", "█   █", " ████", "    █", " ███ "},
	':': {"     ", "  █  ", "     ", "  █  ", "     "},
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	duration, err := parseArgs(args)
	if err != nil {
		return err
	}

	countdown(duration, out)
	return nil
}

func parseArgs(args []string) (time.Duration, error) {
	if len(args) != 1 {
		return 0, errors.New("expected exactly one duration")
	}

	duration, err := time.ParseDuration(args[0])
	if err != nil {
		return 0, fmt.Errorf("invalid duration %q", args[0])
	}
	if duration <= 0 {
		return 0, errors.New("duration must be greater than zero")
	}

	return duration, nil
}

func countdown(duration time.Duration, out io.Writer) {
	defer fmt.Fprint(out, "\x1b[?25h")
	fmt.Fprint(out, "\x1b[?25l")

	deadline := time.Now().Add(duration)
	for {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			renderCountdown(0, out)
			fmt.Fprintln(out, "\nTime's up!")
			return
		}

		renderCountdown(remaining, out)

		sleepFor := time.Second
		if remaining < sleepFor {
			sleepFor = remaining
		}
		time.Sleep(sleepFor)
	}
}

func renderCountdown(duration time.Duration, out io.Writer) {
	fmt.Fprint(out, "\x1b[2J\x1b[H")
	fmt.Fprint(out, "\x1b[38;5;208m")
	fmt.Fprint(out, renderLarge(formatCompact(duration)))
	fmt.Fprint(out, "\x1b[0m")
}

func formatCompact(duration time.Duration) string {
	seconds := int64(duration.Round(time.Second) / time.Second)
	if seconds < 0 {
		seconds = 0
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	remainingSeconds := seconds % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, remainingSeconds)
	}

	return fmt.Sprintf("%02d:%02d", minutes, remainingSeconds)
}

func renderLarge(value string) string {
	var lines [5]string
	for _, char := range value {
		pattern, ok := digitPatterns[char]
		if !ok {
			continue
		}

		for i, row := range pattern {
			if lines[i] != "" {
				lines[i] += "  "
			}
			lines[i] += row
		}
	}

	return strings.Join(lines[:], "\n") + "\n"
}
