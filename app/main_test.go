package main

import (
	"fmt"
	"testing"
)

func ExampleMainHello() {
	fmt.Println("Hello")
	// Output: Hello
}

func TestMainSum(t *testing.T) {
	if 1+2 != 3 {
		t.Fatal("1 + 2 should be 3, but doesn't match")
	}
}
