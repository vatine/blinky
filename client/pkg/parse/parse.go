package parse

import (
	"sort"
)

// Parse a string representation of one or more integers.
// Return a slice of them as int32, in sorted order.
// Uses the following syntax:
//  <n> - number
//  <n>,<expr> - number, plus expression
//  <n>-<m> - all integers between n and m
func Parse(s string) []int32 {
	rv := parseInternal(s, 0, []int32{})
	sort.Slice(rv, func(i, j int) bool { return rv[i] < rv[j] })
	return rv
}

func parseInternal(s string, offset int, acc []int32) []int32 {
	for offset < len(s) {
		n, next := readInt(s, offset)
		if next == len(s) {
			acc = append(acc, n)
			return acc
		}
		if s[next] == ',' {
			offset = next + 1
			acc = append(acc, n)
		}
		if s[next] == '-' {
			m, newNext := readInt(s, next+1)
			more := makeRange(n, m)
			acc = append(acc, more...)
			offset = newNext
		}
	}

	return acc
}

func readInt(s string, offset int) (int32, int) {
	var rv int32

	for ix := offset; ix < len(s); ix++ {
		if s[ix] >= '0' && s[ix] <= '9' {
			rv = 10*rv + int32(s[ix]-'0')
		} else {
			return rv, ix
		}
	}

	return rv, len(s)
}

func makeRange(low, high int32) []int32 {
	rv := make([]int32, 0, high-low+1)

	for n := low; n <= high; n++ {
		rv = append(rv, n)
	}

	return rv
}
