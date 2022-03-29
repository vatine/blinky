package parse

import (
	"testing"
)

func eqSlice(a, b []int32, t *testing.T) bool {
	if len(a) != len(b) {
		t.Errorf("Saw len %d, want len %d", len(a), len(b))
		return false
	}

	for ix, v := range a {
		if v != b[ix] {
			t.Errorf("Position %d, saw %d, want %d", ix, v, b[ix])
			return false
		}
	}

	return true
}

func TestParse(t *testing.T) {
	cases := []struct {
		s    string
		want []int32
	}{
		{"1", []int32{1}},
		{"10", []int32{10}},
		{"1,2", []int32{1, 2}},
		{"2,1", []int32{1, 2}},
		{"1-5", []int32{1, 2, 3, 4, 5}},
	}

	for ix, c := range cases {
		saw := Parse(c.s)
		want := c.want

		if !eqSlice(saw, want, t) {
			t.Errorf("Case %d, saw %v, want %v", ix, saw, want)
		}
	}
}
