// 导入主函数的包
package main

// 系统会导入所需要的包 format 标准输入输出格式包
import "fmt"

// 程序的主入口 程序有且只有一个主函数
func main() {
	// var PI float64 = 3.14159
	var name, age = "张三", 18 //类型推导
	PI := 3.14159            //简短变量声明,只能在函数中使用
	var r float64 = 3
	//var s float64 //声明
	//计算面积
	s := PI * r * r

	fmt.Println("面积", s)
	fmt.Printf("名字：%s \t 年龄:%d", name, age)

}
