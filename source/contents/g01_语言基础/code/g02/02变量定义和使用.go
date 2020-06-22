// 导入主函数的包
package main

// 系统会导入所需要的包 format 标准输入输出格式包
import "fmt"

// 程序的主入口 程序有且只有一个主函数
func main()  {
	// var PI float64 = 3.14159
	PI := 3.14159 //自动推导类型
	var r float64 = 3
	//var s float64 //声明
	//计算面积
	s := PI*r*r 

	fmt.Println("面积",s)
	
}

