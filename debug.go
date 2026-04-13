package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"
)

func main() {
	url := "https://i.oneme.ru/i?r=BTE2sh_eZW7g8kugOdIm2NothKk6IdOFUDZZbQaCXT9rUjAHvcLs7PBUER-zagsNp2s"
	client := &http.Client{Timeout: 10 * time.Second}
	fmt.Printf("🌐 Запрос к: %s\n\n", url)
	
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("❌ СЕТЕВАЯ ОШИБКА: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("⏱ Время ответа: %v\n", duration)
	fmt.Printf("📊 Статус: %s\n\n", resp.Status)

	dump, _ := httputil.DumpResponse(resp, false)
	fmt.Println("--- ЗАГОЛОВКИ ОТВЕТА ---")
	fmt.Println(string(dump))
}
