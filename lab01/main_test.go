package main

import (
	"fmt"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

// func Example() {
// 	main()
// 	// Output:
// 	// this will print to console
// }

// func Example_report() {
// 	report("scruple")
// 	// output:
// 	// U+2108	℈	SCRUPLE
// 	// 1 character found
// }

func Test_CharName_String(t *testing.T) {
	want := "U+0041\tA\tLATIN CAPITAL LETTER A"
	cn := CharName{'A', "LATIN CAPITAL LETTER A"}
	got := fmt.Sprint(cn)
	if got != want {
		t.Errorf("CharName_String\n\tgot:  %q\n\twant: %q", got, want)
	}
}

func Test_contains(t *testing.T) {
	testCases := []struct {
		input []string
		word  string
		want  bool
	}{
		{[]string{"ABC", "DEF"}, "DEF", true},
		{[]string{"ABC", "DEF"}, "GH", false},
	}

	for _, tc := range testCases {
		t.Run(tc.word, func(t *testing.T) {
			got := contains(tc.input, tc.word)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_containsAll(t *testing.T) {
	testCases := []struct {
		input []string
		words []string
		want  bool
	}{
		{[]string{"ABC", "DEF"}, []string{"ABC", "DEF"}, true},
		{[]string{"ABC", "DEF"}, []string{"DEF", "GHI"}, false},
	}

	for _, tc := range testCases {
		t.Run(strings.Join(tc.words, "+"), func(t *testing.T) {
			got := containsAll(tc.input, tc.words)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_scan(t *testing.T) {
	testCases := []struct {
		label string
		start rune
		end   rune
		want  []CharName
	}{
		{"A", 'A', 'B', []CharName{{'A', "LATIN CAPITAL LETTER A"}}},
		{"ABC", 'A', 'D', []CharName{
			{'A', "LATIN CAPITAL LETTER A"},
			{'B', "LATIN CAPITAL LETTER B"},
			{'C', "LATIN CAPITAL LETTER C"},
		}},
		{"Unassigned", '\u0377', '\u037B', []CharName{
			{'\u0377', "GREEK SMALL LETTER PAMPHYLIAN DIGAMMA"},
			{'\u037A', "GREEK YPOGEGRAMMENI"},
		}},
		{"Unnamed", '\x1E', '\x22', []CharName{
			{'\u0020', "SPACE"},
			{'\u0021', "EXCLAMATION MARK"},
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			got := scan(tc.start, tc.end)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_filter(t *testing.T) {
	testCases := []struct {
		start rune
		end   rune
		query []string
		want  []CharName
	}{
		{' ', unicode.MaxRune, []string{"MADEUPWORD"}, []CharName{}},
		{'\u2108', unicode.MaxRune, []string{"SCRUPLE"}, []CharName{{'\u2108', "SCRUPLE"}}},
		{'6', '9', []string{"SEVEN"}, []CharName{{'7', "DIGIT SEVEN"}}},
		{'A', 'C', []string{"A"}, []CharName{{'A', "LATIN CAPITAL LETTER A"}}},
		{',', '/', []string{"MINUS"}, []CharName{{'-', "HYPHEN-MINUS"}}},
		{'?', 'C', []string{"LATIN", "capital"}, []CharName{{'A', "LATIN CAPITAL LETTER A"}, {'B', "LATIN CAPITAL LETTER B"}}},
	}

	for _, tc := range testCases {
		t.Run(strings.Join(tc.query, "+"), func(t *testing.T) {
			sample := scan(tc.start, tc.end)
			got := filter(sample, tc.query)
			assert.Equal(t, tc.want, got)
		})
	}
}
