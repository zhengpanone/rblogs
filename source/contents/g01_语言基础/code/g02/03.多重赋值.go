package main

import "fmt"

func main0301() {
	a,b,c,d := 10,20,30,40
	fmt.Println(a,b,c,d)
}

func main() {
	a,b := 10,20
	fmt.Println(a,b)
	fmt.Println("===========")
	//数据交换
	b,a = a,b
	fmt.Println(a,b)
}