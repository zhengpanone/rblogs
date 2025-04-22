package main

import "fmt"

func add(a int) func(int) int {
	return func(b int) int {
		return a + b
	}
}
func main() {
	result := add(1)
	fmt.Println(result(2))
}
