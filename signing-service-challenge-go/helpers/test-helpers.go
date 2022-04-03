package helpers

import (
	"fmt"
	"testing"
)

func ShouldBe(t *testing.T, value any, expected any) {
	if value != expected {
		t.Fatal(fmt.Sprintf("It was %v when expected %v", value, expected))
	}
}

func ShouldNotBe(t *testing.T, value any, expected any) {
	if value == expected {
		t.Fatal(fmt.Sprintf("It was %v when not expected %v", value, expected))
	}
}
