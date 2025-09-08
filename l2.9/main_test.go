package main

import "testing"

type response struct {
	result string
	err    error
}

type testcase struct {
	input  string
	output response
}

func buildTests() []testcase {
	tescases := []testcase{
		{"a4bc2d5e", response{"aaaabccddddde", nil}},
		{"abcd", response{"abcd", nil}},
		{"45", response{"", ErrNoSymbolBeforeDigit}},
		{"", response{"", nil}},
		{"qwe\\4\\5", response{"qwe45", nil}},
		{"qwe\\45", response{"qwe44444", nil}},
		{"qwe\\", response{"", ErrNoSymbolAfterSlash}},
	}
	return tescases
}

func TestUnpack(t *testing.T) {
	tests := buildTests()
	for _, v := range tests {
		t.Run("testUnpack", func(t *testing.T) {
			data := v
			t.Parallel()
			res, err := unpack(data.input)
			if err != data.output.err || res != data.output.result {
				t.Errorf("input: %s, output: %s, %v, return: %s, %v", data.input, data.output.result, data.output.err, res, err)
			}
		})
	}
}
