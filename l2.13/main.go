package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

const (
	exitErrOccurred = 1
)

var (
	errBadFlagD               = errors.New("the delimiter must be a single character")
	errNoFlagF                = errors.New("you must specify fields")
	errInvalidFieldValue      = errors.New("invalid field value")
	errInvalidDecreasingRange = errors.New("invalid decreasing range")
	errInvalidFieldRange      = errors.New("invalid field range")
	errFieldsNumberedFrom     = errors.New("fields are numbered from 1")
	errInvalidRangeEndpoint   = errors.New("invalid range with no endpoint")
)

type arguments struct {
	f []int
	d string
	s bool
}

func getArgs() (*arguments, error) {
	args := &arguments{}
	pflag.BoolVarP(&args.s, "separated", "s", false, "only strings with delimeter")
	rawD := pflag.StringP("delimeter", "d", "\t", "another delimeter")
	rawF := pflag.StringP("fields", "f", "", "fields to print")

	pflag.Parse()

	if len(*rawD) > 1 {
		return args, errBadFlagD
	}
	args.d = *rawD

	if len(*rawF) == 0 {
		return args, errNoFlagF
	}
	fields := strings.Split(*rawF, ",")
	for i := range fields {
		splitted := strings.Split(fields[i], "-")
		if len(splitted) == 1 {
			conv, err := strconv.Atoi(splitted[0])
			if err != nil {
				return args, fmt.Errorf("%v %s", errInvalidFieldValue, splitted[0])
			}
			if conv < 1 {
				return args, errFieldsNumberedFrom
			}
			args.f = append(args.f, conv)
		} else if len(splitted) == 2 {
			first, err := strconv.Atoi(splitted[0])
			if err != nil {
				return args, errInvalidFieldValue
			}
			second, err := strconv.Atoi(splitted[1])
			if err != nil {
				return args, errInvalidFieldValue
			}
			if second < first {
				return args, errInvalidDecreasingRange
			}
			if first < 1 {
				return args, errFieldsNumberedFrom
			}
			for i := first; i <= second; i++ {
				args.f = append(args.f, i)
			}
		} else {
			return args, errInvalidFieldRange
		}
	}
	slices.Sort(args.f)
	args.f = slices.Compact(args.f)

	return args, nil
}

func cutStrings(data []string, args *arguments) ([]string, error) {
	res := make([]string, 0, len(data))
	for i := range data {
		splitted := strings.Split(data[i], args.d)
		if (args.s && len(splitted) > 1) || (!args.s) {
			toJoin := make([]string, 0, len(args.f))
			for _, v := range args.f {
				if len(splitted) < v {
					break
				}
				toJoin = append(toJoin, splitted[v-1])
			}
			res = append(res, strings.Join(toJoin, args.d))
		}
	}
	return res, nil
}

func main() {
	args, err := getArgs()
	if err != nil {
		fmt.Printf("cut: %v\n", err)
		os.Exit(exitErrOccurred)
	}

	data := make([]string, 0)
	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		data = append(data, sc.Text())
	}

	if sc.Err() != nil {
		fmt.Printf("cut: %v\n", sc.Err())
		os.Exit(exitErrOccurred)
	}

	res, err := cutStrings(data, args)
	if err != nil {
		fmt.Printf("cut: %v\n", err)
	}

	for i := range res {
		fmt.Println(res[i])
	}
}
