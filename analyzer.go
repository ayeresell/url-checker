package main

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func main() {
	url1 := "BTE2sh_eZW7g8kugOdIm2NothKk6IdOFUDZZbQaCXT9rUjAHvcLs7PBUER-zagsNp2s"
	url2 := "BTE2sh_eZW7g8kugOdIm2NotFj26DqUUH5xXRljbhK-CIzAHvcLs7PBUER-zagsNp2s"

	fmt.Printf("Len 1: %d, Len 2: %d\n", len(url1), len(url2))

	d1, _ := decode(url1)
	d2, _ := decode(url2)

	fmt.Println("HEX Comparison:")
	maxLen := len(d1)
	if len(d2) > maxLen { maxLen = len(d2) }

	for i := 0; i < maxLen; i++ {
		b1 := byte(0)
		if i < len(d1) { b1 = d1[i] }
		b2 := byte(0)
		if i < len(d2) { b2 = d2[i] }

		diff := " "
		if b1 != b2 { diff = "!" }
		fmt.Printf("[%02d] %s %02x | %02x\n", i, diff, b1, b2)
	}
}

func decode(s string) ([]byte, error) {
	// Стандартный Base64 URL Safe требует паддинг или Raw
	s = strings.ReplaceAll(s, "-", "+")
	s = strings.ReplaceAll(s, "_", "/")
	
	// Пробуем добавить паддинг если нужно
	for len(s)%4 != 0 {
		s += "="
	}
	return base64.StdEncoding.DecodeString(s)
}
