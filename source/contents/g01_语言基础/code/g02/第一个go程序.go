// 导入主函数的包
package main

//系统导入所需要的包format 标注输入输出格式包
import "fmt"


// 行注释

/*
块注释
*/

var b int // 声明
// 定义变量
var a int = 10

// 程序主入口,程序有且只有一个主函数
func main()  {
	fmt.Print("Hello World!")
	fmt.Println("======================")
	fmt.Println(a)
	fmt.Println(b)
}