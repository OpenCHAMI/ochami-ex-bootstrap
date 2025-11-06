package xname

import "testing"

func TestBMCXnameToNode(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"x1000c0s0b0", "x1000c0s0n0"},
		{"x1000c0s0b1", "x1000c0s0n1"},
		{"x9999c1s2", "x9999c1s2-n0"},
	}
	for _, c := range cases {
		got := BMCXnameToNode(c.in)
		if got != c.out {
			t.Fatalf("BMCXnameToNode(%q)=%q want %q", c.in, got, c.out)
		}
	}
}
