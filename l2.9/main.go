package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrNoSymbolBeforeDigit = fmt.Errorf("ErrNoSymbolBeforeDigit")
	ErrNoSymbolAfterSlash  = fmt.Errorf("ErrNoSymbolAfterSlash")
)

func unpack(s string) (string, error) {
	var result strings.Builder
	data := []rune(s)
	i := 0

	for i < len(data) {
		if unicode.IsDigit(data[i]) {
			if i == 0 || (unicode.IsDigit(data[i-1]) && i > 1 && data[i-2] != '\\') {
				return "", ErrNoSymbolBeforeDigit
			}
			cter, err := strconv.Atoi(string(data[i]))
			if err != nil {
				return "", fmt.Errorf("unpack error: %v", err)
			}
			for range cter - 1 {
				result.WriteRune(data[i-1])
			}
		} else if data[i] == '\\' {
			if i == len(data)-1 {
				return "", ErrNoSymbolAfterSlash
			}
			i += 1
			result.WriteRune(data[i])
		} else {
			result.WriteRune(data[i])
		}
		i += 1
	}
	return result.String(), nil
}

func main() {
	var s string
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	res, err := unpack(s)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Println(res)
}
