<p align="center">
  <img src="https://app-eta-seven-61.vercel.app/banner-urlchecker.svg" width="900"/>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white"/>
  <img src="https://img.shields.io/badge/goroutines-concurrent-00ADD8?style=flat"/>
  <img src="https://img.shields.io/badge/sync%2Fatomic-lock--free-39d353?style=flat"/>
  <img src="https://img.shields.io/badge/Python-3776AB?style=flat&logo=python&logoColor=white"/>
</p>

<h2 align="center">url-checker — High-Performance Concurrent URL Checker</h2>
<p align="center">Takes a base URL, generates thousands of mutations, and checks them concurrently with a goroutine worker pool.</p>

---

## How It Works

```
Base URL
    │
    ▼
Analyzer (mutator)
    ├─ analyzer.go        — character substitution
    ├─ analyzer_bits.go   — bit-level encoding variants
    └─ analyzer_triple.go — triple path segment mutation
    │
    ▼
Worker Pool (goroutines)
    │  ┌──────────────────────────────────┐
    │  │  goroutine 1 ──▶ HTTP GET ──▶ 200? │
    │  │  goroutine 2 ──▶ HTTP GET ──▶ 404? │
    │  │  goroutine N ──▶ HTTP GET ──▶ ...  │
    │  └──────────────────────────────────┘
    │
    ▼
sync/atomic counters (lock-free)
    ├─ totalChecked
    ├─ foundCount
    └─ netErrors
    │
    ▼
found_urls.txt / success.txt
```

## Performance

- Lock-free metrics via `sync/atomic` — no mutex on the hot path
- Connection pooling with custom `http.Transport`
- `sync.Pool` for byte buffers — minimal GC pressure
- Tuner (`tuner.go`) — automatically finds optimal worker count

## Usage

```bash
# Main checker
go build -o checker .
./checker

# Concurrency tuner
go build -o tuner tuner.go
./tuner

# Result parser
python3 parser.py
```

## Project Structure

| File | Description |
|------|-------------|
| `main.go` | Entry point, worker pool, HTTP client |
| `analyzer.go` | URL mutations — character substitution |
| `analyzer_bits.go` | Bit-level and Base64 encoding variants |
| `analyzer_triple.go` | Triple path segment mutation |
| `extractor.go` | Result extraction and filtering |
| `tuner.go` | Concurrency parameter tuner |
| `parser.py` | Log analysis helper |
