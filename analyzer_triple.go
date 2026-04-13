package main

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func main() {
	urls := []string{
		"BTE2sh_eZW7g8kugOdIm2NothKk6IdOFUDZZbQaCXT9rUjAHvcLs7PBUER-zagsNp2s",
		"BTE2sh_eZW7g8kugOdIm2NotFj26DqUUH5xXRljbhK-CIzAHvcLs7PBUER-zagsNp2s",
		"BTE2sh_eZW7g8kugOdIm2NotGUKhTbWRo_fHHQvJonoY_zAHvcLs7PBUER-zagsNp2s",
	}

	var decoded [][]byte
	for _, u := range urls {
		decoded = append(decoded, decode(u))
	}

	fmt.Println("HEX Comparison (3 URLs):")
	for i := 0; i < len(decoded[0]); i++ {
		b1 := decoded[0][i]
		b2 := decoded[1][i]
		b3 := decoded[2][i]

		status := "   "
		if b1 == b2 && b2 == b3 {
			status = "FIX" // Байт совпадает во всех 3 ссылках
		} else {
			status = "VAR" // Байт меняется
		}

		fmt.Printf("[%02d] %s | %02x | %02x | %02x\n", i, status, b1, b2, b3)
	}
}

func decode(s string) []byte {
	s = strings.ReplaceAll(s, "-", "+")
	s = strings.ReplaceAll(s, "_", "/")
	for len(s)%4 != 0 { s += "=" }
	b, _ := base64.StdEncoding.DecodeString(s)
	return b
}
