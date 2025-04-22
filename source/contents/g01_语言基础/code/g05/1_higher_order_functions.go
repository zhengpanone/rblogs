package main

import "fmt"

func filter(numbers []int, callback func(int) bool) []int {

	var result []int
	for _, number := range numbers {
		if callback(number) {
			result = append(result, number)
		}
	}
	return result
}
func isOdd(number int) bool {
	return !isEven(number)
}

func isEven(number int) bool {
	return number%2 == 0
}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	odds := filter(numbers, isOdd)
	fmt.Println(odds)
}
