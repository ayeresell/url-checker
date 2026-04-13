<p align="center">
  <img src="https://app-eta-seven-61.vercel.app/banner-urlchecker.svg" width="900"/>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white"/>
  <img src="https://img.shields.io/badge/goroutines-concurrent-00ADD8?style=flat"/>
  <img src="https://img.shields.io/badge/sync%2Fatomic-lock--free-39d353?style=flat"/>
  <img src="https://img.shields.io/badge/Python-3776AB?style=flat&logo=python&logoColor=white"/>
</p>

<h2 align="center">url-checker — высокопроизводительный URL-чекер на Go</h2>
<p align="center">Инструмент для конкурентной мутации и проверки URL. Берёт базовую ссылку, генерирует тысячи вариантов и проверяет их параллельно.</p>

---

## Как это работает

```
Базовый URL
    │
    ▼
Analyzer (мутатор)
    ├─ analyzer.go       — подстановка символов
    ├─ analyzer_bits.go  — bit-level кодировки
    └─ analyzer_triple.go — тройная мутация сегментов
    │
    ▼
Worker Pool (горутины)
    │  ┌─────────────────────────────────┐
    │  │  goroutine 1 ──▶ HTTP GET ──▶ 200? │
    │  │  goroutine 2 ──▶ HTTP GET ──▶ 404? │
    │  │  goroutine N ──▶ HTTP GET ──▶ ...  │
    │  └─────────────────────────────────┘
    │
    ▼
sync/atomic счётчики (без мьютексов)
    ├─ totalChecked
    ├─ foundCount
    └─ netErrors
    │
    ▼
found_urls.txt / success.txt
```

## Производительность

- Lock-free метрики через `sync/atomic` — нет блокировок на горячем пути
- Пул соединений через кастомный `http.Transport`
- `sync.Pool` для буферов байт — минимальная нагрузка на GC
- Тюнер (`tuner.go`) — автоматический подбор оптимального числа воркеров

## Использование

```bash
# Основной чекер
go build -o checker .
./checker

# Тюнер параметров конкурентности
go build -o tuner tuner.go
./tuner

# Парсер результатов
python3 parser.py
```

## Структура проекта

| Файл | Описание |
|------|----------|
| `main.go` | Точка входа, worker pool, HTTP-клиент |
| `analyzer.go` | Мутации URL — посимвольная замена |
| `analyzer_bits.go` | Bit-level и Base64 варианты кодировки |
| `analyzer_triple.go` | Тройная мутация сегментов пути |
| `extractor.go` | Извлечение и фильтрация результатов |
| `tuner.go` | Подбор оптимальной конкурентности |
| `parser.py` | Анализ логов |
