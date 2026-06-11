# Timer

[![CI](https://github.com/ragnarok22/timer-go/actions/workflows/ci.yml/badge.svg)](https://github.com/ragnarok22/timer-go/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/ragnarok22/timer-go/graph/badge.svg)](https://codecov.io/gh/ragnarok22/timer-go)
![Go](https://img.shields.io/badge/go-1.22+-00ADD8?logo=go&logoColor=white)
![CLI](https://img.shields.io/badge/type-CLI-orange)

A small terminal countdown timer with a large seven-segment-style display.

## Install

With Go:

```sh
go install github.com/ragnarok22/timer-go@latest
```

This installs the CLI as `timer-go`.

With the install script:

```sh
curl -fsSL https://raw.githubusercontent.com/ragnarok22/timer-go/main/install.sh | sh
```

This installs the CLI as `timer` from the latest GitHub release.

## Usage

```sh
timer <duration>
```

Examples:

```sh
timer 10s
timer 6m
timer 6m30s
timer 1h
```

Durations use Go's duration format, so combinations like `1h5m10s` also work.

## Display

The timer uses a compact format:

- Under 1 hour: `MM:SS`
- 1 hour or more: `H:MM:SS`

## Help

```sh
timer help
timer -h
timer --help
```

## Cancel

Press `Ctrl+C` at any time to cancel the timer cleanly.

## Development

Run tests:

```sh
go test ./...
```

Run tests with coverage:

```sh
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out
```
