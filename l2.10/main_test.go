package main

import (
	"testing"

	"slices"
)

type input struct {
	data []string
	args *arguments
}

type response struct {
	data []string
	err  error
}

type testcase struct {
	name     string
	in       input
	expected response
}

func buildDefaultTestcases() []testcase {
	return []testcase{
		{name: "4fun", in: input{[]string{"ba"}, &arguments{1, false, false, false}}, expected: response{[]string{"ba"}, nil}},
		{name: "just strings", in: input{[]string{"ba bb", "a", "7", "ba aa", "ba"}, &arguments{1, false, false, false}}, expected: response{[]string{"7", "a", "ba", "ba aa", "ba bb"}, nil}},
		{name: "-u -k flags", in: input{[]string{"a b", "c a", "a b", "b c", "b c"}, &arguments{k: 2, n: false, r: false, u: true}}, expected: response{[]string{"c a", "a b", "b c"}, nil}},
		{name: "+ -r flag", in: input{[]string{"a b", "c a", "a b", "b c", "b c"}, &arguments{2, false, true, true}}, expected: response{[]string{"b c", "a b", "c a"}, nil}},
		{name: "-n flag", in: input{[]string{"a b", "a a", "a", "b", "30", "4", "3 b", "3 a", "3"}, &arguments{1, true, false, false}}, expected: response{[]string{"a", "a a", "a b", "b", "3", "3 a", "3 b", "4", "30"}, nil}},
	}
}

func TestSortStringsWrap(t *testing.T) {
	defaultTestcases := buildDefaultTestcases()
	t.Run("default", func(t *testing.T) {
		for i := range defaultTestcases {
			t.Run(defaultTestcases[i].name, func(t *testing.T) {
				old := slices.Clone(defaultTestcases[i].in.data)
				out, err := sortStringsWrap(defaultTestcases[i].in.data, defaultTestcases[i].in.args)
				if !slices.Equal(out, defaultTestcases[i].expected.data) || err != defaultTestcases[i].expected.err {
					t.Errorf("%v\n%v\n%v", old, defaultTestcases[i].expected.data, out)
				}
			})
		}
	})

}
