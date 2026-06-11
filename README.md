# Timer

[![CI](https://github.com/ragnarok22/timer-go/actions/workflows/ci.yml/badge.svg)](https://github.com/ragnarok22/timer-go/actions/workflows/ci.yml)
![Go](https://img.shields.io/badge/go-1.22+-00ADD8?logo=go&logoColor=white)
![CLI](https://img.shields.io/badge/type-CLI-orange)

A small terminal countdown timer with a large seven-segment-style display.

## Install

```sh
go install .
```

This installs the CLI as `timer`.

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
