package main

import "fmt"

func main1() {
	var a int
	//& 运算符,取地址运算符
	fmt.Scan(&a)
	//%p 占位符,表示输出一个数据对应的内存地址
	// ox表示十六进制数据
	fmt.Println(a)
	fmt.Printf("%p", &a)

	var s1,s2 string

	fmt.Scan(&s1, &s2)
	fmt.Println(s1+"\n"+s2)
}


func main2(){
	var r float64
	PI := 3.14159
	fmt.Printf("请输入半径:")
	fmt.Scanf("%f",&r)
	fmt.Println(r)
	fmt.Printf("面积：%2.f\n",PI*r*r)
}

func main(){
	main2()
}