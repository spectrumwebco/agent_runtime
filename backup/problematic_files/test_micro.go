package main

import (
	"fmt"
)

func main() {
	fmt.Println("Testing Go Micro integration")
	r := registry.NewRegistry()
	fmt.Printf("Registry initialized: %v\n", r != nil)
}
