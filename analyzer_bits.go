package main

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func main() {
	// Декодируем ваши две рабочие ссылки
	s1 := "BTE2sh_eZW7g8kugOdIm2NothKk6IdOFUDZZbQaCXT9rUjAHvcLs7PBUER-zagsNp2s"
	s2 := "BTE2sh_eZW7g8kugOdIm2NotFj26DqUUH5xXRljbhK-CIzAHvcLs7PBUER-zagsNp2s"

	d1 := decode(s1)
	d2 := decode(s2)

	fmt.Println("HEX Comparison of Payload (Bytes 18-33):")
	for i := 18; i <= 33; i++ {
		b1 := d1[i]
		b2 := d2[i]
		
		// Побитовый XOR чтобы увидеть разницу
		xor := b1 ^ b2
		
		fmt.Printf("[%02d] S1: %08b (%02x) | S2: %08b (%02x) | XOR: %08b\n", i, b1, b1, b2, b2, xor)
	}
}

func decode(s string) []byte {
	s = strings.ReplaceAll(s, "-", "+")
	s = strings.ReplaceAll(s, "_", "/")
	for len(s)%4 != 0 { s += "=" }
	b, _ := base64.StdEncoding.DecodeString(s)
	return b
}
