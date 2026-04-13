# url-checker

High-performance concurrent URL checker written in Go. Generates and tests URL variants to find valid endpoints.

## What it does

- Takes a base URL and mutates it systematically (character substitution, encoding variants)
- Sends concurrent HTTP requests with connection pooling
- Logs found valid URLs and tracks success/error rates in real time
- Includes a tuner mode for adjusting concurrency parameters

## Tech Stack

- **Go** — concurrent requests via goroutines
- **sync/atomic** — lock-free counters for performance metrics
- **Python** — helper parser script for result analysis

## Usage

```bash
go build -o checker .
./checker
```

Tuner mode:
```bash
go build -o tuner tuner.go
./tuner
```
