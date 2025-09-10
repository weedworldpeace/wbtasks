package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

const (
	exitErrorOccurred = 1
	defaultColumn     = 1
	spaceSymbol       = " "
)

type arguments struct {
	k int
	n bool
	r bool
	u bool
}

func (a *arguments) sanitize() {
	if a.k < 1 { // column must be natural
		a.k = defaultColumn
	}
}

func columnStrings(data []string, args *arguments) []string { // func to put privilege column to start of a slice
	front := args.k
	if front > len(data) {
		front = defaultColumn
	}
	return slices.Concat(data[front-1:], data[:front-1])
}

func uncolumnStrings(data []string, args *arguments) []string {
	front := args.k
	if front > len(data) {
		front = defaultColumn
	}
	return slices.Concat(data[len(data)-(front-1):], data[:len(data)-(front-1)])
}

func cmpTwoStrings(a, b string) int {
	front := len(a)
	if len(b) < front {
		front = len(b)
	}
	for i := range front {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		} else {
			return 0
		}
	}
	if len(a) < len(b) {
		return -1
	} else if len(b) < len(a) {
		return 1
	} else {
		return 0
	}
}

func cmpPrivilege(a, b []string) int {
	front := len(a)
	if len(b) < front {
		front = len(b)
	}
	for i := range front {
		if i == 0 { // only privilege column may be sorted like nums
			first, ferr := strconv.Atoi(a[i])
			second, serr := strconv.Atoi(b[i])
			if ferr != nil && serr != nil { // two strings
				res := cmpTwoStrings(a[i], b[i])
				if res == 1 || res == -1 {
					return res
				}
			} else if ferr != nil && serr == nil { // string and int
				return -1
			} else if serr != nil && ferr == nil { // int and string
				return 1
			} else {
				if first < second { // two int
					return -1
				} else if first > second {
					return 1
				}
			}
		} else {
			res := cmpTwoStrings(a[i], b[i])
			if res == 1 || res == -1 {
				return res
			}
		}
	}
	if len(a) < len(b) {
		return -1
	} else if len(b) < len(a) {
		return 1
	} else {
		return 0
	}
}

func sortStringsLikeNums(data []string, args *arguments) ([]string, error) {
	dataToSort := make([][]string, len(data))
	for i := range data {
		dataToSort[i] = columnStrings(strings.Split(data[i], spaceSymbol), args)
	}

	slices.SortFunc(dataToSort, cmpPrivilege)

	for i := range data {
		data[i] = strings.Join(dataToSort[i], spaceSymbol)
		data[i] = strings.Join(uncolumnStrings(strings.Split(data[i], spaceSymbol), args), spaceSymbol)
	}
	return data, nil
}

func sortStrings(data []string, args *arguments) ([]string, error) {
	for i := range data {
		data[i] = strings.Join(columnStrings(strings.Split(data[i], spaceSymbol), args), spaceSymbol)
	}
	slices.Sort(data)
	for i := range data {
		data[i] = strings.Join(uncolumnStrings(strings.Split(data[i], spaceSymbol), args), spaceSymbol)
	}
	return data, nil
}

func uniqueStrings(data []string) []string {
	hm := make(map[string]bool)
	idx := 0
	for i := range data {
		if _, b := hm[data[i]]; !b {
			data[idx] = data[i]
			hm[data[i]] = true
			idx += 1
		}
	}
	return data[:idx:idx]
}

func sortStringsWrap(data []string, args *arguments) ([]string, error) {
	var result []string
	var err error
	if args.u {
		data = uniqueStrings(data)
	}
	if args.n {
		result, err = sortStringsLikeNums(data, args)
	} else {
		result, err = sortStrings(data, args)
	}
	if args.r {
		slices.Reverse(result)
	}
	return result, err
}

func getArgs() *arguments {
	args := arguments{}
	flag.IntVarP(&args.k, "kolumn", "k", defaultColumn, "number of column to sort")
	flag.BoolVarP(&args.n, "numeric", "n", false, "sort like numbers")
	flag.BoolVarP(&args.r, "reverse", "r", false, "reverse sort")
	flag.BoolVarP(&args.u, "unique", "u", false, "leave only unique")

	flag.Parse()

	return &args
}

func main() {
	data := make([]string, 0)
	args := getArgs()
	args.sanitize()

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		data = append(data, sc.Text())
	}

	res, err := sortStringsWrap(data, args)
	if err != nil {
		fmt.Println(err)
		os.Exit(exitErrorOccurred)
	}
	for _, v := range res {
		fmt.Println(v)
	}
}
