package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

const (
	baseURL        = "https://i.oneme.ru/i?r="
	prefix         = "BTE2sh_eZW7g8kugOdIm2Not"
	suffix         = "AHvcLs7PBUER-zagsNp2s"
	errorThreshold = 0.05
	testDuration   = 5
)

type Stats struct {
	Total        uint64
	Errors       uint64
	LatencyNanos uint64
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("🚀 ЗАПУСК УМНОГО ТЮНЕРА СЕТИ")
	fmt.Println("--------------------------------------------------")

	targetRPS := 100
	workers := 20
	bestRPS, bestWorkers := 0, 0

	for {
		fmt.Printf("📊 ТЕСТ: Цель %d RPS | Потоков: %d\n", targetRPS, workers)
		stats := runTestStep(targetRPS, workers)
		actualRPS := int(stats.Total / testDuration)
		errRate := 0.0
		if stats.Total > 0 {
			errRate = float64(stats.Errors) / float64(stats.Total)
		}
		if stats.Total == 0 { break }

		fmt.Printf("   📈 Факт: %d RPS | Ошибки: %.1f%%", actualRPS, errRate*100)
		if errRate > errorThreshold {
			fmt.Printf("\n🛑 ПРЕДЕЛ СЕРВЕРА: на %d RPS пошли блокировки.\n", targetRPS)
			break
		}

		bestRPS, bestWorkers = targetRPS, workers
		efficiency := float64(actualRPS) / float64(targetRPS)
		if efficiency < 0.95 {
			workers = int(float64(workers)/efficiency) + 15
			fmt.Printf(" -> 🆙 Добавляем воркеры -> %d\n", workers)
		} else {
			fmt.Printf(" -> ✅ Ок, повышаем нагрузку...\n")
			targetRPS += 100
			if workers < targetRPS/2 { workers += 20 }
		}
		time.Sleep(1 * time.Second)
		fmt.Println("--------------------------------------------------")
	}

	fmt.Println("\n🏆 ИДЕАЛЬНЫЕ НАСТРОЙКИ:")
	fmt.Printf("rpsLimit   = %d\nnumWorkers = %d\n", bestRPS, bestWorkers)
}

func runTestStep(targetRPS int, workerCount int) Stats {
	var s Stats
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(testDuration)*time.Second)
	defer cancel()

	transport := &http.Transport{
		DialContext: (&net.Dialer{Timeout: 2 * time.Second}).DialContext,
		MaxIdleConns: workerCount, MaxIdleConnsPerHost: workerCount,
	}
	client := &http.Client{Transport: transport, Timeout: 4 * time.Second}
	jobs := make(chan string, targetRPS*2)
	
	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done(): return
				case url := <-jobs:
					atomic.AddUint64(&s.Total, 1)
					req, _ := http.NewRequestWithContext(ctx, "HEAD", url, nil)
					req.Header.Set("User-Agent", "Mozilla/5.0")
					start := time.Now()
					resp, err := client.Do(req)
					if err != nil { atomic.AddUint64(&s.Errors, 1) } else {
						resp.Body.Close()
						atomic.AddUint64(&s.LatencyNanos, uint64(time.Since(start).Nanoseconds()))
					}
				}
			}
		}()
	}

	ticker := time.NewTicker(time.Second / time.Duration(targetRPS))
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done(): return s
		case <-ticker.C:
			select {
			case jobs <- generateURL():
			default:
			}
		}
	}
}

func generateURL() string {
	b := make([]byte, 16)
	rand.Read(b)
	mid := base64.RawURLEncoding.EncodeToString(b)
	if len(mid) > 22 { mid = mid[:22] }
	return baseURL + prefix + mid + suffix
}
