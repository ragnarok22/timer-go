# AGENTS.md

## Project Shape
- Single-package Go CLI (`module timer`, `package main`) targeting Go 1.22.
- Runtime entrypoint is `main.go`; tests live in `main_test.go` and exercise unexported helpers directly.
- No external dependencies are currently declared; there is no `go.sum` unless dependencies are added.

## Commands
- No Makefile or task runner is present; use Go commands directly.
- Install locally: `go install .`
- Run a short manual timer: `go run . 1s`
- Run all tests: `go test ./...`
- Run one test: `go test . -run TestFormatCompact`
- CI-equivalent checks, in order: `go mod download`, `test -z "$(gofmt -l .)"`, `go vet ./...`, `go test ./...`
- Format changed Go files with `gofmt -w <files>`; CI only checks `gofmt`, not `go fmt` output.

## Implementation Notes
- Duration parsing intentionally uses `time.ParseDuration`; bare numbers like `10` are invalid, while `10s`, `6m30s`, and `1h` are valid.
- Display format is `MM:SS` under one hour and `H:MM:SS` at one hour or more.
- Terminal rendering writes ANSI escape codes for clear-screen, color, and cursor visibility; tests assert exact byte substrings, so preserve restore behavior on cancellation and completion.
- `countdown` accepts an interrupt channel for testability; avoid replacing it with direct signal handling inside the loop.

## CI
- GitHub Actions workflow is `.github/workflows/ci.yml` and runs on push and pull request.
- Setup uses `actions/setup-go` with `go-version-file: go.mod`, so update `go.mod` when changing the Go version.
