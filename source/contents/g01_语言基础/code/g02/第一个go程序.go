// 包, 表明代码所在的模块
package main

//系统导入所需要的包format 标注输入输出格式包
import "fmt"

// 行注释

/*
块注释
*/

// 函数外面只能放标识符(变量、常量、函数、类型)的声明,不能放语句

var b int // 声明
// 定义变量
var a int = 10

// 批量声明
var (
	name string
	age  int
	isOK bool
)

// 程序主入口,程序有且只有一个主函数
func main() {
	// GO语言变量声明必须使用,非全局变量声明了必须使用
	name = "张三"
	age = 15
	isOK = true
	fmt.Print("Hello World!")
	fmt.Println("======================")
	fmt.Println(a)
	fmt.Println(b)

	fmt.Println(name)

}
