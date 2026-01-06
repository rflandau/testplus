package testplus_test

import (
	"testing"
	"testplus"
)

func TestPortGenerator_Generate(t *testing.T) {
	// test that we panic on precisely the (65535-1023=)64512nd call
	pg := testplus.NewPortGenerator()
	seen := make(map[uint16]struct{}, 64512)
	for range 64512 {
		n := pg.Generate()
		if _, found := seen[n]; found {
			t.Fatal("repeat digit found: ", n)
		}
		seen[n] = struct{}{}
	}
	// with all numbers exhausted, we should now see a panic
	var noPanic bool // will not be set if we panicked as expected
	defer func() {
		recover()
		if noPanic {
			t.Fatal("did not panic on final Generate call")
		}
	}()
	pg.Generate()
	noPanic = true

}

func TestSlicesUnorderedEqual(t *testing.T) {
	tests := []struct {
		name  string
		a     []any
		b     []any
		equal bool
	}{
		{"arrays are already equal", []any{"Hello", "World"}, []any{"Hello", "World"}, true},
		{"arrays are inverted", []any{"Hello", "World"}, []any{"World", "Hello"}, true},
		{"arrays are unequal", []any{"Hello"}, []any{"World"}, false},
		{"different lengths", []any{"Hello", "World"}, []any{"World"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testplus.SlicesUnorderedEqual(tt.a, tt.b); got != tt.equal {
				t.Errorf("SlicesUnorderedEqual() = %v, want %v", got, tt.equal)
			}
		})
	}
}
