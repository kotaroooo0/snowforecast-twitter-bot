package text

import (
	"fmt"
	"testing"
)

func ExampleTextHello() {
	fmt.Println("Hello")
	// Output: Hello
}

func TestTextSum(t *testing.T) {
	if 1+2 != 3 {
		t.Fatal("1 + 2 should be 3, but doesn't match")
	}
}
