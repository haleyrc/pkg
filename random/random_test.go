package random_test

import (
	"testing"

	"github.com/haleyrc/cheevos/internal/lib/random"
)

func TestStringReturnsRandomStrings(t *testing.T) {
	var (
		length = 32
		first  = random.String(length)
		second = random.String(length)
	)

	if len(first) != length {
		t.Fatalf("Expected String to return a %d character string, but got %d character string %q.", length, len(first), first)
	}

	if len(second) != length {
		t.Fatalf("Expected String to return a %d character string, but got %d character string %q.", length, len(second), second)
	}

	if first == second {
		t.Errorf("Expected String to return unique strings, but it didn't.")
	}
}
