package handlers

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"Motif Sasirangan", "motif-sasirangan"},
		{"Naga Balimbur!", "naga-balimbur"},
		{"  Quiz 101  ", "quiz-101"},
		{"", ""},
	}
	for _, tc := range tests {
		if got := slugify(tc.in); got != tc.want {
			t.Fatalf("slugify(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
