<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/concurrent-goroutines-00ADD8?style=flat" />
  <img src="https://img.shields.io/badge/Python-3776AB?style=flat&logo=python&logoColor=white" />
</p>

<h1 align="center">url-checker</h1>
<p align="center">High-performance concurrent URL checker in Go. Mutates a base URL across thousands of variants and finds valid endpoints in seconds.</p>

---

## How It Works

1. Takes a known valid URL as a seed
2. Generates variants by mutating path segments (character substitution, Base64 encoding variants, bit flips)
3. Fires concurrent HTTP requests using a goroutine worker pool
4. Tracks hits, misses, and network errors in real time with atomic counters
5. Writes successful URLs to output files

## Performance

- Lock-free metrics via `sync/atomic`
- Connection pooling with custom `http.Transport`
- `sync.Pool` for byte buffer reuse — zero GC pressure at high concurrency

## Usage

**Checker:**
```bash
go build -o checker .
./checker
```

**Tuner** (concurrency parameter search):
```bash
go build -o tuner tuner.go
./tuner
```

**Result parser:**
```bash
python3 parser.py
```

## Project Structure

| File | Description |
|------|-------------|
| `main.go` | Entry point, worker pool, HTTP client |
| `analyzer.go` | URL mutation logic |
| `analyzer_bits.go` | Bit-level encoding variants |
| `analyzer_triple.go` | Triple-segment mutation |
| `extractor.go` | Result extraction helpers |
| `tuner.go` | Concurrency parameter tuner |
| `parser.py` | Python log parser |
