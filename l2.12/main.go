package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/pflag"
)

var (
	errNoRegexp = errors.New("no regexp")
)

const (
	exitErrOccurred = 1
	occured         = "error occured"
)

type arguments struct {
	regexpMask string
	A          int
	B          int
	C          int
	c          bool
	i          bool
	v          bool
	F          bool
	n          bool
}

type response struct {
	rawNumber int
	val       string
}

func (a *arguments) valid() error {
	if a.A == 0 {
		a.A = a.C
	}
	if a.B == 0 {
		a.B = a.C
	}
	return nil
}

func getArgs() (*arguments, error) {
	args := &arguments{}

	pflag.IntVarP(&args.A, "A", "A", 0, "strings after")
	pflag.IntVarP(&args.B, "B", "B", 0, "strings before")
	pflag.IntVarP(&args.C, "C", "C", 0, "strings before and after")
	pflag.BoolVarP(&args.c, "c", "c", false, "only quantity of result")
	pflag.BoolVarP(&args.i, "i", "i", false, "ignore case")
	pflag.BoolVarP(&args.v, "v", "v", false, "invert filter")
	pflag.BoolVarP(&args.F, "F", "F", false, "fixed string")
	pflag.BoolVarP(&args.n, "n", "n", false, "print with row number")

	pflag.Parse()

	if len(pflag.Args()) != 1 {
		return args, errNoRegexp
	}
	args.regexpMask = pflag.Args()[0]

	err := args.valid()

	return args, err
}

func remakeDataContext(args *arguments, i, splittedSize int) (int, int) {
	rawsBefore := args.B
	rawsAfter := args.A
	if i < args.B {
		rawsBefore = i
	}
	if splittedSize-i-1 < args.A {
		rawsAfter = splittedSize - i - 1
	}
	return rawsBefore, rawsAfter
}

func compactString(splitted []string, i, rawsBefore, rawsAfter int) string {
	var bd strings.Builder
	for j := range rawsBefore {
		bd.WriteString(splitted[i-rawsBefore+j])
		bd.WriteByte('\n')
	}
	bd.WriteString(splitted[i])
	for j := range rawsAfter {
		bd.WriteByte('\n')
		bd.WriteString(splitted[i+j+1])
	}
	return bd.String()
}

func grep(data string, args *arguments) ([]response, error) {
	res := make([]response, 0)
	splitted := strings.Split(data, "\n")
	if args.F {
		args.regexpMask = regexp.QuoteMeta(args.regexpMask)
	}
	if args.i {
		args.regexpMask = fmt.Sprintf("(?i)%s", args.regexpMask)
	}
	re, err := regexp.Compile(args.regexpMask)
	if err != nil {
		return []response{}, fmt.Errorf("grep :%v", err)
	}
	for i, v := range splitted {
		rawsBefore, rawsAfter := remakeDataContext(args, i, len(splitted))
		matched := re.MatchString(v)
		if (args.v && !matched) || (!args.v && matched) {
			compacted := compactString(splitted, i, rawsBefore, rawsAfter)
			res = append(res, response{i + 1, compacted})
		}
	}
	return res, nil
}

func writeResult(res []response, args *arguments) {
	if args.c {
		fmt.Println(len(res))
	} else {
		for _, v := range res {
			if args.n {
				fmt.Printf("%d:%s\n", v.rawNumber, v.val)
			} else {
				fmt.Printf("%s\n", v.val)
			}
		}
	}
}

func main() {
	args, err := getArgs()
	if err != nil {
		fmt.Printf("%s: %v\n", occured, err)
		os.Exit(exitErrOccurred)
	}
	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("%s: %v\n", occured, err)
		os.Exit(exitErrOccurred)
	}

	res, err := grep(string(raw), args)
	if err != nil {
		fmt.Printf("%s: %v\n", occured, err)
		os.Exit(exitErrOccurred)
	}
	writeResult(res, args)
}
