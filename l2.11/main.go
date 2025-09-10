package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

var (
	errBadSymbol = errors.New("bad symbol error")
)

const (
	rusAlphabetSize = 33
)

type el struct {
	first  string
	values []string
}

func findAnagrams(data []string, dict map[rune]int) []el {
	hm := make(map[[rusAlphabetSize]int]el)
	for i := range data {
		freqs := [rusAlphabetSize]int{}
		for _, v := range data[i] {
			freqs[dict[v]] += 1
		}
		if _, b := hm[freqs]; b {
			hm[freqs] = el{hm[freqs].first, append(hm[freqs].values, data[i])}
		} else {
			hm[freqs] = el{first: data[i], values: []string{data[i]}}
		}
	}
	res := make([]el, 0, len(hm))
	for _, v := range hm {
		if len(v.values) > 1 {
			slices.Sort(v.values)
			res = append(res, v)
		}
	}
	return res
}

func buildDict() map[rune]int {
	dict := make(map[rune]int)
	start := 'а'
	for i := range rusAlphabetSize {
		dict[start] = i
		start += 1
	}
	return dict
}

func valid(s string, dict map[rune]int) error {
	for _, v := range s {
		if _, b := dict[v]; !b {
			return errBadSymbol
		}
	}
	return nil
}

func main() {
	dict := buildDict()

	data := make([]string, 0)
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		raw := strings.ToLower(strings.TrimSpace(sc.Text()))
		err := valid(raw, dict)
		if err == nil {
			data = append(data, raw)
		}
	}

	res := findAnagrams(data, dict)

	for _, v1 := range res {
		fmt.Printf("%s: [ ", v1.first)
		for _, v2 := range v1.values {
			fmt.Printf("%s ", v2)
		}
		fmt.Printf("]\n")
	}
}
