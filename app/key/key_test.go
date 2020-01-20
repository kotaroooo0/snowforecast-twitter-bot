package key

import (
	"fmt"
	"testing"
)

func ExampleKeyHello() {
	fmt.Println("Hello")
	// Output: Hello
}

func TestKeySum(t *testing.T) {
	if 1+2 != 3 {
		t.Fatal("1 + 2 should be 3, but doesn't match")
	}
}
