package main

import "fmt"

// 定义一个算子
func Map[T any, R any](list []T, callback func(T) R) []R {
	result := make([]R, len(list))
	for index, number := range list {
		result[index] = callback(number)
	}
	return result
}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// 定义一个算子
	even := Map(numbers, func(number int) int {
		return number * 2
	})
	fmt.Println(even)
}
