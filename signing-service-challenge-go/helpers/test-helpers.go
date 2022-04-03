package helpers

import (
	"fmt"
	"testing"
)

func ShouldBe(t *testing.T, value any, expected any) {
	if value != expected {
		t.Fatal(fmt.Sprintf("It was %s when expected %s", value, expected))
	}
}

func ShouldNotBe(t *testing.T, value any, expected any) {
	if value == expected {
		t.Fatal(fmt.Sprintf("It was %s when expected %s", value, expected))
	}
}
